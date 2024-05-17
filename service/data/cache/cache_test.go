package cache

// Shouldn't have key 1
// func TestBase(t *testing.T) {
// 	cac := NewCache(100)
// 	for i := 0; i < 1000; i++ {
// 		order := &models.Order{Id: strconv.Itoa(i)}
// 		cac.Store(strconv.Itoa(i), order)
// 	}
// 	if _, err := cac.Get("1"); err != NotExistError {
// 		t.Failed()
// 	}
// 	if _, err := cac.Get("950"); err == NotExistError {
// 		t.Failed()
// 	}
// }
