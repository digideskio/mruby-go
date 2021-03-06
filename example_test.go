// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

package mruby_test

import (
	"fmt"

	"github.com/olivere/mruby-go"
)

func Example() {
	// Create a new context
	ctx := mruby.NewContext()

	// Run a script that returns the value of 1+2
	val, err := ctx.LoadString("1 + 2")
	if err != nil {
		fmt.Println("LoadStringResult failed")
		return
	}

	// Convert the value to an int
	i, err := val.ToInt()
	if err != nil {
		fmt.Println("Result is not a Fixnum")
		return
	}
	fmt.Printf("1+2=%d\n", i)
	// Output: 1+2=3
}

func Example_automaticallyConvertToGoInterface() {
	// Create a new context
	ctx := mruby.NewContext()

	// Run a script that returns the value of 1+2 and directly converts to Go
	res, err := ctx.LoadStringResult("1 + 2")
	if err != nil {
		fmt.Println("LoadStringResult failed")
		return
	}
	fmt.Printf("1+2=%d", res)
	// Output: 1+2=3
}

func Example_withArguments() {
	// Create a new context.
	ctx := mruby.NewContext()

	// Run a script that adds all values and directly converts the result to Go.
	// The args of LoadStringXXX will be available via ARGV in the script.
	res, err := ctx.LoadStringResult("ARGV.inject { |x,y| x+y }", 1, 2, 3.5)
	if err != nil {
		fmt.Println("LoadStringResult failed")
		return
	}
	fmt.Println(res)
	// Output: 6.5
}

func ExampleContext_New() {
	// Create a new context, and set some options
	ctx := mruby.NewContext(mruby.SetFilename("test.rb"), mruby.SetNoExec(true))
	if ctx != nil {
		fmt.Println("Context initialized")
	}
	// Output: Context initialized
}

func ExampleParser() {
	// Create a new context
	ctx := mruby.NewContext()

	// Create a parser to parse the given code
	code := `
def concat(a, b)
	a + b
end

concat "Hello", "World"
`
	parser, err := ctx.Parse(code)
	if err != nil {
		fmt.Println("Parser cannot parse")
		return
	}

	// Run the code
	val, err := parser.Run()
	if err != nil {
		fmt.Println("Run failed")
		return
	}
	res, err := val.ToInterface()
	if err != nil {
		fmt.Println("Cannot convert")
		return
	}
	s, ok := res.(string)
	if !ok {
		fmt.Println("Not a string")
		return
	}
	fmt.Println(s)
	// Output: HelloWorld
}

func ExampleFunction() {
	// Create a new context, and set some options
	ctx := mruby.NewContext()
	if ctx == nil {
		fmt.Println("Cannot initialize context")
		return
	}

	// sayHello is an extension method that can be called from Ruby.
	sayHello := func(ctx *mruby.Context) (output mruby.Value, err error) {
		// We expect a string here.
		args, err := ctx.GetArgs()
		if err != nil {
			return mruby.NilValue(ctx), err
		}
		if len(args) == 1 {
			s, err := args[0].ToString()
			if err != nil {
				return mruby.NilValue(ctx), err
			}
			s = fmt.Sprintf("Hello %s!", s)
			return ctx.ToValue(s)
		}
		return mruby.NilValue(ctx), nil
	}

	// We create a new module called Helpers that will hold our extension method.
	module, err := ctx.DefineModule("Helpers", nil)
	if err != nil {
		fmt.Println("Cannot create module")
		return
	}
	if module == nil {
		fmt.Println("Method is nil")
		return
	}

	// Now we register the extension method. It will be available as
	// Helper.say_hello from within Ruby and requires 1 argument.
	module.DefineClassMethod("say_hello", sayHello)

	greeting, err := ctx.LoadStringResult("Helpers.say_hello(ARGV[0])", "Matz")
	if err != nil {
		fmt.Printf("Cannot execute say_hello method: %v\n", err)
		return
	}
	fmt.Println(greeting)
	// Output: Hello Matz!
}
