package services

import "fmt"

type Avamar struct {
	cmd *Command
	prefix string
}

func NewAvamar(cmd *Command) *Avamar {
	return &Avamar{
		cmd: cmd,
		prefix: "[Response]:",
	}
}

// avmaint sched stop --ava
// avmaint checkpoint --ava --wait
// avmaint --ava garbagecollect
// avmaint --ava gcstatus
func (a *Avamar) Gc() (error) {
	output, err := a.cmd.Run("avmaint sched stop --ava")
	if err != nil {
		fmt.Printf("Avamar Gc Error: %s\n", err)
		return err
	}
	fmt.Println(a.prefix, output)
	output, err = a.cmd.Run("avmaint checkpoint --ava --wait")
	fmt.Println(a.prefix, output)
	if err != nil {

		fmt.Printf("Avamar Gc Error: %s\n", err)
		return err
	}
	output, err = a.cmd.Run("avmaint --ava garbagecollect")
	fmt.Println(a.prefix, output)
	if err != nil {
		fmt.Printf("Avamar Gc Error: %s\n", err)
		return err
	}
	output, err = a.cmd.Run("avmaint --ava gcstatus")
	fmt.Println(a.prefix, output)
	if err != nil {
		fmt.Printf("Avamar Gc Error: %s\n", err)
		return err
	}
	return nil
}