package snweb

import (
	"log"
	"reflect"
	"testing"
)

func TestRoute(t *testing.T) {

	expectedHandler1 := func(context *HttpContext) { log.Panicln("1") }
	expectedRefHandler1 := reflect.ValueOf(expectedHandler1)

	expectedHandler2 := func(context *HttpContext) { log.Panicln("2") }
	expectedRefHandler2 := reflect.ValueOf(expectedHandler2)

	rh := RoutedHandler{}
	rh.AddRoute("/", expectedHandler1)
	rh.AddRoute("/project/:projectId", expectedHandler2)

	route1 := rh.FindMatchingRoute("/")
	reflectHandler1 := reflect.ValueOf(route1.Handler)
	if reflectHandler1.Pointer() != expectedRefHandler1.Pointer() {
		t.Error("Can't find proper function for URL /")
	}

	route2 := rh.FindMatchingRoute("/project/1")
	reflectHandler2 := reflect.ValueOf(route2.Handler)
	if reflectHandler2.Pointer() != expectedRefHandler2.Pointer() {
		t.Error("Can't find proper function for URL /project/1")
	}
	projectId := route2.PathParams["projectId"]
	if projectId != "1" {
		t.Error("Path param should be '1'")
	}

}

