package services

import (
	"datadomain/pkg/parser"
	"fmt"
)

type DataDomain struct {
	cmd    *Command
	prefix string
}

func NewDataDomain(cmd *Command) *DataDomain {
	return &DataDomain{
		cmd:    cmd,
		prefix: "[Response]:",
	}
}

func (d *DataDomain) FilesysGc() error {
	output, err := d.cmd.Run("filesys clean start")
	fmt.Println(d.prefix, output)
	if err != nil {
		fmt.Printf("DataDomain FilesysGc Error: %s\n", err)
		return err
	}
	output, err = d.cmd.Run("filesys clean watch")
	fmt.Println(d.prefix, output)
	if err != nil {
		fmt.Printf("DataDomain FilesysGc Error: %s\n", err)
		return err
	}
	return nil
}

func (d *DataDomain) DataMoveToCloud() error {
	output, err := d.cmd.Run("data-movement start to-tier cloud")
	fmt.Println(d.prefix, output)
	if err != nil {
		fmt.Printf("DataDomain DataMoveToCloud Error: %s\n", err)
		return err
	}
	output, err = d.cmd.Run("data-movement watch")
	fmt.Println(d.prefix, output)
	if err != nil {
		fmt.Printf("DataDomain DataMoveToCloud Error: %s\n", err)
		return err
	}
	return nil
}

func (d *DataDomain) DDBoostRestart() error {
	output, err := d.cmd.Run("ddboost disable")
	fmt.Println(d.prefix, output)
	if err != nil {
		fmt.Printf("DataDomain DDBoostRestart Error: %s\n", err)
		return err
	}
	output, err = d.cmd.Run("ddboost enable")
	fmt.Println(d.prefix, output)
	if err != nil {
		fmt.Printf("DataDomain DDBoostRestart Error: %s\n", err)
		return err
	}
	return nil
}

func (d *DataDomain) DeleteAllStorageUnit() error {
	output, err := d.cmd.Run("ddboost storage-unit show")
	fmt.Println(d.prefix, output)
	if err != nil {
		fmt.Printf("DataDomain DeleteAllStorageUnit: %s\n", err)
		return err
	}

	length, storages := parser.DDBoostStorageUnitShow(output)
	fmt.Printf("%s Total: %d storage(s)\n", d.prefix, length)
	for _, storage := range storages {
		d.cmd.Run("ddboost storage-unit delete %s", storage)
	}

	err = d.FilesysGc()
	if err != nil {
		fmt.Printf("DataDomain DeleteAllStorageUnit: %s\n", err)
		return err
	}

	return nil
}

func (d *DataDomain) DeleteOldestStorageUnit() error {
	output, err := d.cmd.Run("ddboost storage-unit show")
	fmt.Println(d.prefix, output)
	if err != nil {
		fmt.Printf("DataDomain DeleteAllStorageUnit: %s\n", err)
		return err
	}

	length, storages := parser.DDBoostStorageUnitShow(output)
	fmt.Printf("%sTotal: %d storages(s)\n", d.prefix, length)
	for _, storage := range storages {
		d.cmd.Run("ddboost storage-unit delete %s", storage)
		break
	}

	err = d.FilesysGc()
	if err != nil {
		fmt.Printf("DataDomain DeleteAllStorageUnit: %s\n", err)
		return err
	}

	return nil
}
