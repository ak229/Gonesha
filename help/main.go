package help

func Needed(args []string) bool{
 
	if(len(args) == 0 || args[0] == "--help") {
         	
		return true
 	
	} else {

		return false
	
	}
}
func Greet() {


}
