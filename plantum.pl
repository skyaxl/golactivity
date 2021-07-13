@startuml
start
partition Execute {
	:array1 := [][]int({INT(1)});
	:i := INT(0);
	repeat
		:array2 := array1[i];
		:j := len(array2);
		if (j == INT(0)) then (yes)
			end
			note right:Return (fmt.Sprint(j))
		else
		endif
	repeat while (i < len(array1)) is (true)
	:b := true;
	if ((request.OK == b || ctx == nil) && request.Name != "nullable") then (yes)
		:fs.ExecuteNone(ctx,request);
		end
		note right:Return ("no")
	else
	endif
	:c1 := make(chan string);
	:validator := func () ( bool){Literal func};
	:validator();

	switch (switch(request.Name))
	case (case "a")
		end
		note right:Return (request.Name + request.Name)

	case (case "go")
		end
		note right:Return (<-c1)

	case (default)
		end
		note right:Return (request.Name)

	endswitch
}
stop
@enduml