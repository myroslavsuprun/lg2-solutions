package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type bufferWg struct {
	buf *bytes.Buffer
	wg  *sync.WaitGroup
}

func (bWg *bufferWg) Write(p []byte) (int, error) {
	defer bWg.wg.Done()
	return bWg.buf.Write(p)
}

func TestWriteData(t *testing.T) {
	wg := sync.WaitGroup{}

	b := &bufferWg{
		wg:  &wg,
		buf: bytes.NewBuffer([]byte{}),
	}
	ch := make(chan Result, 100)

	go WriteData(ch, b)

	input := []Result{
		{"11", 22},
		{"iasdkhajfhdasd", math.MaxInt},
		{"4ef5e2ab-900d-46e8-8eae-ef55fd8c440b", 128318238},
		{"12832183", -28312831823},
		{"iajskhgAS@@#!5/32@#%@#($#)@!", 123123},
	}

	want := bytes.NewBuffer([]byte{})
	for _, i := range input {
		ch <- i
		wg.Add(1)

		_, err := want.WriteString(fmt.Sprintf("%s:%d\n", i.Id, i.Value))
		if err != nil {
			t.Errorf("got error: %v", err)
		}
	}
	close(ch)

	wg.Wait()

	if !bytes.Equal(b.buf.Bytes(), want.Bytes()) {
		t.Errorf("expected: %v, got: %v", want, b.buf)
	}
}

func TestParser(t *testing.T) {
	data := []struct {
		name string
		in   []byte
		exp  Input
		err  bool
	}{
		{"1", []byte("someId\n+\n299\n22"), Input{"someId", "+", 299, 22}, false},
		{"2", []byte("sss\n-\n88\n12"), Input{"sss", "-", 88, 12}, false},
		{"3", []byte("4ef5e2ab-900d-46e8-8eae-ef55fd8c440b\n+\n9999999999999\n1"), Input{"4ef5e2ab-900d-46e8-8eae-ef55fd8c440b", "+", 9999999999999, 1}, false},
		{"4", []byte("sss\n/\n88\n4"), Input{"sss", "/", 88, 4}, false},
		{"5", []byte("8542jjdjsdf\n*\n22\n8"), Input{"8542jjdjsdf", "*", 22, 8}, false},
		{"6", []byte("i52\n85235\n12\n12"), Input{}, true},
		{"7", []byte("i53\n-\n88\nj"), Input{}, true},
		{"8", []byte("i54\n-\nn234\n22"), Input{}, true},
		{"9", []byte("i55\n-\n88\n12\n-"), Input{}, true},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			input, err := parser(d.in)

			if err != nil && !d.err {
				t.Errorf("expected: no error, got: %v", err)
				return
			}

			if diff := cmp.Diff(d.exp, input); diff != "" {
				t.Errorf("expected: %v, got: %v", d.exp, input)
				return
			}
		})

	}
}

func TestDataProcessor_Conc(t *testing.T) {
	data := getTestData()

	in := make(chan []byte, 100)
	out := make(chan Result, 100)

	go DataProcessor(in, out)

	expRes := []Result{}
	gotRes := []Result{}

	for _, d := range data {
		in <- d.in

		if !d.skip {
			expRes = append(expRes, d.res)
		}
	}

	close(in)
	for {
		v, ok := <-out
		if !ok {
			break
		}
		gotRes = append(gotRes, v)
	}

	if dif := cmp.Diff(expRes, gotRes); dif != "" {
		t.Errorf("expected: %v, got: %v", expRes, gotRes)
	}
}

func getTestData() []struct {
	name string
	in   []byte
	res  Result
	skip bool
} {
	return []struct {
		name string
		in   []byte
		res  Result
		skip bool
	}{
		{"1+1=2", []byte("1\n+\n1\n1"), Result{"1", 2}, false},
		{"2*2=4", []byte("2\n*\n2\n2"), Result{"2", 4}, false},
		{"6/2=3", []byte("3\n/\n6\n2"), Result{"3", 3}, false},
		{"6-2=4", []byte("4\n-\n6\n2"), Result{"4", 4}, false},
		{"122+001=123", []byte("5\n+\n122\n001"), Result{"5", 123}, false},
		{"22*2=44", []byte("6\n*\n22\n2"), Result{"6", 44}, false},
		{"66/2=33", []byte("7\n/\n66\n2"), Result{"7", 33}, false},
		{"6-22=-16", []byte("8\n-\n6\n22"), Result{"8", -16}, false},
		{"4-x", []byte("9\n-\n4\nx"), Result{}, true},
		{"4%4", []byte("10\n%\n4\n4"), Result{}, true},
	}

}

func TestDataProcessor_Seq(t *testing.T) {
	data := getTestData()

	in := make(chan []byte, 100)
	out := make(chan Result, 100)

	go DataProcessor(in, out)

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			in <- d.in

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
			defer cancel()
			select {
			case <-ctx.Done():
				if !d.skip {
					t.Errorf("expected: %v, got: %v", d.res, ctx.Err())
				}
			case res := <-out:
				if d.skip {
					t.Errorf("expected: skip, got: %v", res)
					return
				}

				if dif := cmp.Diff(d.res, res); dif != "" {
					t.Errorf("expected: %v, got: %v", d.res, res)
				}
			}
		})
	}
}

func getNewControllerData() []struct {
	name     string
	body     []byte
	respCode int
	respBody []byte
} {
	return []struct {
		name     string
		body     []byte
		respCode int
		respBody []byte
	}{
		{"empty body", []byte{}, 400, []byte("Bad Input")},
		{"OK1", []byte("123 123 1231 2312 3"), 202, []byte("OK: 2")},
		{"OK2", []byte("632\n+\n77\n23"), 202, []byte("OK: 3")},
		{"OK3", []byte("632\n*\n1723\n23"), 202, []byte("OK: 4")},
		{"OK3", []byte("632\n/\n12312737\n23"), 202, []byte("OK: 5")},
	}
}

func TestNewController_channel(t *testing.T) {
	data := getNewControllerData()

	got := make(chan []byte, len(data))
	want := make(chan []byte, len(data))
	h := NewController(got)
	for _, d := range data {
		if d.respCode >= 200 && d.respCode < 300 {
			want <- d.body
		}

		_, err := makeNewControllerReq(h, d.body)
		if err != nil {
			t.Errorf("expected: http request, got: %v", err)
			return
		}
	}
	close(got)
	close(want)

	var gotB []byte
	for i := range got {
		gotB = append(gotB, i...)
	}
	var wantB []byte
	for i := range want {
		wantB = append(wantB, i...)
	}

	if !bytes.Equal(gotB, wantB) {
		t.Errorf("expected: %v, got: %v", wantB, gotB)
	}
}

func TestNewController_backup(t *testing.T) {
	out := make(chan []byte, 100)
	h := NewController(out)

	for i := 0; i < 200; i++ {
		rr, err := makeNewControllerReq(h, []byte("some data"))
		if err != nil {
			t.Errorf("expected: http request, got: %v", err)
			return
		}

		if i == 150 {
			for k := 0; k < 100; k++ {
				<-out
			}
		}

		if i > 99 && i < 150 && rr.Code != 503 {
			t.Errorf("expected: http code 503, got: %v", rr.Code)
			return
		}

		if i > 150 && rr.Code != 202 {
			t.Errorf("expected: http code 202, got: %v", rr.Code)
			return
		}

	}
}

func TestNewController_race(t *testing.T) {
	out := make(chan []byte, 100)
	h := NewController(out)

	for i := 0; i < 200; i++ {
		go func() {
			_, err := makeNewControllerReq(h, []byte("some data"))
			if err != nil {
				t.Errorf("expected: http request, got: %v", err)
				return
			}

		}()
	}
}

func TestNewController_pressure(t *testing.T) {
	out := make(chan []byte, 100)
	h := NewController(out)

	for i := 0; i < 200; i++ {
		rr, err := makeNewControllerReq(h, []byte("some data"))
		if err != nil {
			t.Errorf("expected: http request, got: %v", err)
			return
		}

		if i <= 99 {
			continue
		}
		if rr.Code != 503 {
			t.Errorf("expected: http code 503, got: %v", rr.Code)
			return
		}

		exp := fmt.Sprintf("Too Busy: %d", i-99)
		if rr.Body.String() != exp {
			t.Errorf("expected: %v, got: %v", exp, rr.Code)
			return
		}

	}
}

func TestNewController_response(t *testing.T) {
	data := getNewControllerData()

	out := make(chan []byte, 100)
	h := NewController(out)
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			rr, err := makeNewControllerReq(h, d.body)
			if err != nil {
				t.Errorf("expected: http request, got: %v", err)
				return
			}

			if rr.Code != d.respCode {
				t.Errorf("expected: %v, got: %v", d.respCode, rr.Code)
				return
			}

			if rr.Body.String() != string(d.respBody) {
				t.Errorf("expected: %v, got: %v", d.respBody, rr.Body)
				return
			}
		})
	}

}

func makeNewControllerReq(h http.Handler, rBody []byte) (*httptest.ResponseRecorder,
	error) {
	req, err := http.NewRequest(http.MethodGet, "/", bytes.NewReader(rBody))
	if err != nil {
		return &httptest.ResponseRecorder{}, err
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr, nil
}

func FuzzParser(f *testing.F) {
	testcases := [][]byte{
		[]byte("12\n+\n222\n751725"),
		[]byte("usiofhsdiofhsfohsdofihsfpo2341904324\n*\n4823424\n8419421841294"),
		[]byte("hasd\n%\n833\n123123"),
		[]byte("jssakdjasd\n123812312jksdjasdkasd\nasdjasdkj\njjsdasd"),
		[]byte("jssakdjasd\n1282jksad\nasjdsadjk\ntuaiusaidusaidu"),
	}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		input, err := parser(b)
		if err != nil {
			t.Skip("handled error")
		}

		reB := []byte(fmt.Sprintf("%v\n%v\n%d\n%d", input.Id, input.Op,
			input.Val1, input.Val2))
		newInp, err := parser(reB)
		if err != nil {
			t.Errorf("expected: pass, got: %v", err)
			return
		}

		if diff := cmp.Diff(input, newInp); diff != "" {
			t.Errorf("expected: %v, got: %v", input, newInp)
		}
	})

}
