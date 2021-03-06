package command

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strconv"
	"text/template"
)

func writeFileFromTemplate(tplSrc, dst string, perm os.FileMode,
	data interface{}, chownTo *user.User) {
	b, err := Asset(tplSrc)
	swallow(err)
	if data != nil {
		wr := &bytes.Buffer{}
		t := template.Must(template.New(tplSrc).Parse(string(b)))
		err = t.Execute(wr, data)
		swallow(err)

		if err = ioutil.WriteFile(dst, wr.Bytes(), perm); err != nil {
			os.MkdirAll(path.Dir(dst), 0755)
		}

		err = ioutil.WriteFile(dst, wr.Bytes(), perm)
		swallow(err)

		return
	}

	// no template, just file copy
	if err = ioutil.WriteFile(dst, b, perm); err != nil {
		os.MkdirAll(path.Dir(dst), 0755)
	}

	err = ioutil.WriteFile(dst, b, perm)
	swallow(err)

	if chownTo != nil {
		chown(dst, chownTo)
	}
}

func chown(fp string, chownTo *user.User) {
	uid, _ := strconv.Atoi(chownTo.Uid)
	gid, _ := strconv.Atoi(chownTo.Gid)
	swallow(os.Chown(fp, uid, gid))
}
