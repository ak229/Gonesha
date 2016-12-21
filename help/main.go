package help

func Needed(args *[]string) bool{

	arguments := *args  
	if(len(arguments) == 0 || arguments[0] == "--help") {
         	
		return true
 	
	} else {

		return false
	
	}
}
