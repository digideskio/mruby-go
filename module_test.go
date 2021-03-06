// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

package mruby

import (
	"html"
	"testing"
)

func TestNewModule(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	mod, err := NewModule(ctx, "Helpers", nil)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if mod == nil {
		t.Errorf("expected module; got: %v", mod)
	}

	if !ctx.HasModule("Helpers", nil) {
		t.Fatalf("expected to find module %q", "Helpers")
	}
}

func TestDefineModule(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	mod, err := ctx.DefineModule("Helpers", nil)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if mod == nil {
		t.Errorf("expected module; got: %v", mod)
	}

	if !ctx.HasModule("Helpers", nil) {
		t.Fatalf("expected to find module %q", "Helpers")
	}
}

func TestDefineModuleScoping(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	_, found := ctx.GetModule("MissingModule", nil)
	if found {
		t.Errorf("expected to not find module %q; got: %v", "MissingModule", found)
	}

	outer, err := ctx.DefineModule("Outer", nil)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if outer == nil {
		t.Fatalf("expected outer module; got: %v", outer)
	}
	_, found = ctx.GetModule("Outer", nil)
	if !found {
		t.Errorf("expected to find module %q; got: %v", "Outer", found)
	}
	if !ctx.HasModule("Outer", nil) {
		t.Errorf("expected to find module %q", "Outer")
	}

	inner, err := ctx.DefineModule("Inner", outer)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if inner == nil {
		t.Fatalf("expected inner module; got: %v", inner)
	}
	_, found = ctx.GetModule("Inner", outer)
	if !found {
		t.Errorf("expected to find module %q; got: %v", "Outer::Inner", found)
	}
	_, found = ctx.GetModule("Inner", nil)
	if found {
		t.Errorf("expected to not find module %q; got: %v", "::Inner", found)
	}
	if !ctx.HasModule("Inner", outer) {
		t.Errorf("expected to find module %q", "Outer::Inner")
	}
}

func TestModuleDefineMethodWithNoArgs(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	mod, err := ctx.DefineModule("Helpers", nil)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if mod == nil {
		t.Errorf("expected module; got: %v", mod)
	}

	helloWorld := func(ctx *Context) (Value, error) {
		return ctx.ToValue("Hello world")
	}

	mod.DefineClassMethod("hello", helloWorld)

	s, err := ctx.LoadStringResult("Helpers.hello()")
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello world" {
		t.Errorf("expected %q; got: %q", "Hello world", s)
	}
}

func TestModuleDefineMethodWithRequiredArg(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	mod, err := ctx.DefineModule("Helpers", nil)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if mod == nil {
		t.Errorf("expected module; got: %v", mod)
	}

	escapeHtml := func(ctx *Context) (output Value, err error) {
		// We expect a string here.
		args, err := ctx.GetArgs()
		if err != nil {
			return NilValue(ctx), err
		}
		if len(args) == 1 {
			s, err := args[0].ToString()
			if err == nil {
				s = html.EscapeString(s)
			}
			return ctx.ToValue(s)
		}
		return NilValue(ctx), nil
	}

	mod.DefineClassMethod("escape_html", escapeHtml)

	input := "<esc&ped>"
	expected := html.EscapeString(input)

	got, err := ctx.LoadStringResult("Helpers.escape_html(ARGV[0])", input)
	if err != nil {
		t.Fatal(err)
	}
	if got != expected {
		t.Errorf("expected %q; got: %q", expected, got)
	}
}
