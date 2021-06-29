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
	repeat
		:key := keyOf array1;
		:v := itemOf array1;
		:fmt.Print(key,v);
	repeat while (range array1) is (true)
	:b := true;
	if ((request.OK == b || ctx == nil) && request.Name != STRING("nullable")) then (yes)
		:fs.ExecuteNone(ctx,request);
		end
		note right:Return (STRING("no"))
	else
	endif
	if (INT(2) * INT(2) == INT(4)) then (yes)
		end
		note right:Return (STRING("4"))
	else
	endif
	if (!request.OK) then (yes)
		end
		note right:Return (STRING("none"))
	else
	endif
	if (request.Name == STRING("")) then (yes)
		end
		note right:Return (STRING("name"))
	else
	endif
	:fs.ExecuteNone(ctx,request);
	:validator := ;
	:validator();
	end
	note right:Return (STRING(""))
}
stop
@enduml