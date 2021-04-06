package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
)

func multiplier(ctx context.Context, number int) int {
	span, _ := opentracing.StartSpanFromContext(ctx, "multiplier")
	defer span.Finish()
	time.Sleep(700 * time.Millisecond)
	return number * number

}
func adder(ctx context.Context, number int) int {
	span, _ := opentracing.StartSpanFromContext(ctx, "adder")
	defer span.Finish()
	time.Sleep(400 * time.Millisecond)
	return number + 1
}

func calculate(ctx context.Context, number int) int {
	span, _ := opentracing.StartSpanFromContext(ctx, "calculate")
	defer span.Finish()
	calculateCtx := opentracing.ContextWithSpan(ctx, span)

	time.Sleep(1 * time.Second)

	newNum := multiplier(calculateCtx, number)

	time.Sleep(1 * time.Second)

	newNum = adder(calculateCtx, newNum)

	time.Sleep(1 * time.Second)

	return newNum
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("app")
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(r.Context(), span)

	keys, ok := r.URL.Query()["number"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'number' is missing")
		return
	}
	key := keys[0]
	log.Println("Url Param 'number' is: " + string(key))

	number, err := strconv.ParseInt(key, 10, 0)
	if err != nil {
		fmt.Fprintf(w, "bad parameter")
		return
	}

	response := calculate(ctx, int(number))
	fmt.Fprintf(w, fmt.Sprint(response))
}
