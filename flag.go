// Copyright 2014 The erzha Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package kernel


import (
	"fmt"
	"os"
	"flag"
	"github.com/erzha/econf"
	"github.com/erzha/elog"
)

var flagBasedir *string
var flagConffile *string

//set the main flags for all kinds of servers
func setMainFlags() {
	flagBasedir	 = flag.String("basedir", "", "set the basedir, the $(pwd) is default")
	flagConffile = flag.String("conf", "erzha.ini", "$basedir/conf/erzha.ini is the default.")
}

func parseArgs() {
	if flag.Parsed() {
		return
	}
	
	flag.Usage = printHelpInfo
	flag.Parse()

	if "" == *flagBasedir {
		var err error
		*flagBasedir, err = os.Getwd()
		if nil != err {
			fmt.Println("error occurs when getting pwd:", err)
			os.Exit(-1)
		}
	}

	if "" == *flagConffile {
		*flagConffile = "erzha.ini"
	}
}

func printHelpInfo() {
	fmt.Fprintf(os.Stderr, "erzha is a framework written in go.\n")
	fmt.Fprintf(os.Stderr, "version: %s\n", Version)
	fmt.Fprintf(os.Stderr, "\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
}


func initConf() *econf.Conf {
	var confFile string
	var pConf *econf.Conf
	var err error

	pConf = econf.NewConf()
	confFile = *flagBasedir + "/conf/" + *flagConffile
	err = pConf.ParseFile(confFile)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	return pConf
}

func initLogger(conf *econf.Conf) *elog.Logger {
	logger := elog.NewLogger()
	
	level := conf.String(Conf_key_erzha_log_level, "info")
	logger.SetMinLogLevelByName(level)
	
	outFile := conf.String(Conf_key_erzha_log_file, "")
	writer := os.Stdout
	if "" != outFile {
		f, err := os.Open(outFile)
		if nil != err {
			logger.Fatalf("open_log_file_failed filename:%s err:%s", outFile, err.Error())
		}
		writer = f
	}
	
	logger.SetLogWriter(writer)
	
	return logger
}

func init() {
	setMainFlags()
	}