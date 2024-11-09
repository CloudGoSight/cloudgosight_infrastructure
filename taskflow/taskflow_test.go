package taskflow_test

import (
	"fmt"
	"github.com/CloudGoSight/cloudgosight_infrastructure/taskflow"
	"os"
	"testing"
	"time"
)

var exector = taskflow.NewExecutor(10)

func TestTaskFlow(t *testing.T) {
	A, B, C :=
		taskflow.NewTask("A", func() {
			fmt.Println("A")
		}),
		taskflow.NewTask("B", func() {
			fmt.Println("B")
		}),
		taskflow.NewTask("C", func() {
			fmt.Println("C")
		})

	A1, B1, C1 :=
		taskflow.NewTask("A1", func() {
			fmt.Println("A1")
		}),
		taskflow.NewTask("B1", func() {
			fmt.Println("B1")
		}),
		taskflow.NewTask("C1", func() {
			fmt.Println("C1")
		})
	A.Precede(B)
	C.Precede(B)
	A1.Precede(B)
	C.Succeed(A1)
	C.Succeed(B1)

	tf := taskflow.NewTaskFlow("G")
	tf.Push(A, B, C)
	tf.Push(A1, B1, C1)

	exector.Run(tf).Wait()
	fmt.Print("########### second times")
	exector.Run(tf).Wait()
}

func TestSubflow(t *testing.T) {
	A, B, C :=
		taskflow.NewTask("A", func() {
			fmt.Println("A")
		}),
		taskflow.NewTask("B", func() {
			fmt.Println("B")
		}),
		taskflow.NewTask("C", func() {
			fmt.Println("C")
		})

	A1, B1, C1 :=
		taskflow.NewTask("A1", func() {
			fmt.Println("A1")
		}),
		taskflow.NewTask("B1", func() {
			fmt.Println("B1")
		}),
		taskflow.NewTask("C1", func() {
			fmt.Println("C1")
		})
	A.Precede(B)
	C.Precede(B)
	A1.Precede(B)
	C.Succeed(A1)
	C.Succeed(B1)

	subflow := taskflow.NewSubflow("sub1", func(sf *taskflow.Subflow) {
		A2, B2, C2 :=
			taskflow.NewTask("A2", func() {
				fmt.Println("A2")
			}),
			taskflow.NewTask("B2", func() {
				fmt.Println("B2")
			}),
			taskflow.NewTask("C2", func() {
				fmt.Println("C2")
			})
		A2.Precede(B2)
		C2.Precede(B2)
		sf.Push(A2, B2, C2)
	})

	subflow2 := taskflow.NewSubflow("sub2", func(sf *taskflow.Subflow) {
		A3, B3, C3 :=
			taskflow.NewTask("A3", func() {
				fmt.Println("A3")
			}),
			taskflow.NewTask("B3", func() {
				fmt.Println("B3")
			}),
			taskflow.NewTask("C3", func() {
				fmt.Println("C3")
				// time.Sleep(10 * time.Second)
			})
		A3.Precede(B3)
		C3.Precede(B3)
		sf.Push(A3, B3, C3)
	})

	subflow.Precede(B)
	subflow.Precede(subflow2)

	tf := taskflow.NewTaskFlow("G")
	tf.Push(A, B, C)
	tf.Push(A1, B1, C1, subflow, subflow2)
	exector.Run(tf)
	exector.Wait()
	//if err := taskflow.Visualize(tf, os.Stdout); err != nil {
	//	log.Fatal(err)
	//}
	exector.Profile(os.Stdout)
	// exector.Wait()

	// if err := tf.Visualize(os.Stdout); err != nil {
	// 	panic(err)
	// }
}

// ERROR robust testing
func TestTaskflowPanic(t *testing.T) {
	A, B, C :=
		taskflow.NewTask("A", func() {
			fmt.Println("A")
		}),
		taskflow.NewTask("B", func() {
			fmt.Println("B")
		}),
		taskflow.NewTask("C", func() {
			fmt.Println("C")
			panic("panic C")
		})
	A.Precede(B)
	C.Precede(B)
	tf := taskflow.NewTaskFlow("G")
	tf.Push(A, B, C)

	exector.Run(tf).Wait()
}

func TestSubflowPanic(t *testing.T) {
	A, B, C :=
		taskflow.NewTask("A", func() {
			fmt.Println("A")
		}),
		taskflow.NewTask("B", func() {
			fmt.Println("B")
		}),
		taskflow.NewTask("C", func() {
			fmt.Println("C")
		})
	A.Precede(B)
	C.Precede(B)

	subflow := taskflow.NewSubflow("sub1", func(sf *taskflow.Subflow) {
		A2, B2, C2 :=
			taskflow.NewTask("A2", func() {
				fmt.Println("A2")
				time.Sleep(1 * time.Second)
			}),
			taskflow.NewTask("B2", func() {
				fmt.Println("B2")
			}),
			taskflow.NewTask("C2", func() {
				fmt.Println("C2")
				panic("C2 paniced")
			})
		sf.Push(A2, B2, C2)
		A2.Precede(B2)
		panic("subflow panic")
		C2.Precede(B2)
	})

	subflow.Precede(B)

	tf := taskflow.NewTaskFlow("G")
	tf.Push(A, B, C)
	tf.Push(subflow)
	exector.Run(tf)
	exector.Wait()
	//if err := taskflow.Visualize(tf, os.Stdout); err != nil {
	//	fmt.Errorf("%v", err)
	//}
	exector.Profile(os.Stdout)
}

func TestTaskflowCondition(t *testing.T) {
	A, B, C :=
		taskflow.NewTask("A", func() {
			fmt.Println("A")
		}),
		taskflow.NewTask("B", func() {
			fmt.Println("B")
		}),
		taskflow.NewTask("C", func() {
			fmt.Println("C")
		})
	A.Precede(B)
	C.Precede(B)
	tf := taskflow.NewTaskFlow("G")
	tf.Push(A, B, C)
	fail, success := taskflow.NewTask("failed", func() {
		fmt.Println("Failed")
		t.Fail()
	}), taskflow.NewTask("success", func() {
		fmt.Println("success")
	})

	cond := taskflow.NewCondition("cond", func() uint { return 0 })
	B.Precede(cond)
	cond.Precede(success, fail)

	suc := taskflow.NewSubflow("sub1", func(sf *taskflow.Subflow) {
		A2, B2, C2 :=
			taskflow.NewTask("A2", func() {
				fmt.Println("A2")
			}),
			taskflow.NewTask("B2", func() {
				fmt.Println("B2")
			}),
			taskflow.NewTask("C2", func() {
				fmt.Println("C2")
			})
		sf.Push(A2, B2, C2)
		A2.Precede(B2)
		C2.Precede(B2)
	})
	fs := taskflow.NewTask("fail_single", func() {
		fmt.Println("it should be canceled")
	})
	fail.Precede(fs, suc)
	// success.Precede(suc)
	tf.Push(cond, success, fail, fs, suc)
	exector.Run(tf).Wait()

	//if err := taskflow.Visualize(tf, os.Stdout); err != nil {
	//	fmt.Errorf("%v", err)
	//}
	exector.Profile(os.Stdout)
}

func TestTaskflowLoop(t *testing.T) {
	A, B, C :=
		taskflow.NewTask("A", func() {
			fmt.Println("A")
		}),
		taskflow.NewTask("B", func() {
			fmt.Println("B")
		}),
		taskflow.NewTask("C", func() {
			fmt.Println("C")
		})
	A.Precede(B)
	C.Precede(B)
	tf := taskflow.NewTaskFlow("G")
	tf.Push(A, B, C)
	zero := taskflow.NewTask("zero", func() {
		fmt.Println("zero")
	})
	counter := uint(0)
	cond := taskflow.NewCondition("cond", func() uint {
		counter += 1
		return counter % 3
	})
	B.Precede(cond)
	cond.Precede(cond, cond, zero)

	tf.Push(cond, zero)
	exector.Run(tf).Wait()

	//if err := taskflow.Visualize(tf, os.Stdout); err != nil {
	//	fmt.Errorf("%v", err)
	//}
	exector.Profile(os.Stdout)
}
