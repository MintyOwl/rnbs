package main

import (
	"bytes"
	"io/ioutil"
	"strings"
)

type replace struct {
	old    string
	new    string
	source string
	dest   string
	folder bool
}

func newReplace(old, new, source, dest string) *replace {
	return &replace{
		old:    old,
		new:    new,
		source: source,
		dest:   dest,
	}
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

func (rep *replace) copy() error {
	oldbyte, oldlower := []byte(trim(rep.old)), []byte(strings.ToLower(trim(rep.old)))
	newbyte, newlower := []byte(trim(rep.new)), []byte(strings.ToLower(trim(rep.new)))
	fl, err := ioutil.ReadFile(rep.source)
	bt := bytes.Replace(fl, oldbyte, newbyte, -1)
	bt = bytes.Replace(bt, oldlower, newlower, -1)
	ioutil.WriteFile(rep.dest, bt, 755)
	return err
}

func folderName(old, new, source, sroot, droot string) string {

	//fmt.Print(ok)
	str := strings.Replace(source, trim(old), trim(new), -1)
	str = strings.Replace(str, strings.ToLower(trim(old)), strings.ToLower(trim(new)), -1)
	str = strings.Replace(str, sroot, droot, -1)
	return str

}
