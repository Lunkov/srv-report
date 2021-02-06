package main

import (
  "os"
  "bytes"
  "fmt"
  "io/ioutil"
  "html/template"
  "github.com/jung-kurt/gofpdf"
  "github.com/golang/glog"
)

var memTemplate = make(map[string]*template.Template)

func makeReport(cfg *ConfigInfo, report_settings *ReportInfo, prop *map[string]string) (map[string]string, bool) {
  if glog.V(2) {
    glog.Infof("LOG: REPORT: run prepare...\n")
  }
  result := make(map[string]string)
  var ok bool
  var okLocal bool
  var okBuffer bool
  var err error
  var repLang string
  var repBuffer string
  var repTemplate string
  var repLocalFile string
  var repFormat string
  var repTitle string
   
  repTemplate, ok = (*prop)["REPORT_TEMPLATE"]
  if !ok {
    glog.Errorf("ERR: REPORT: 'REPORT_TEMPLATE' don`t set\n")
    return result, false
  }
  
  repLang, ok = (*prop)["REPORT_LANG"]
  if !ok {
    glog.Warningf("ERR: REPORT: 'REPORT_LANG' don`t set\n")
    repLang = cfg.DefaultLang
  }

  repLocalFile, okLocal = (*prop)["REPORT_RESULT_LOCAL_FILE"]
  repBuffer, okBuffer = (*prop)["REPORT_RESULT_BUFFER"]
  
  if !okLocal && !okBuffer {
    glog.Errorf("ERR: REPORT: 'REPORT_RESULT_LOCAL_FILE' don`t set\n")
    return result, false
  }

  repFormat, ok = (*prop)["REPORT_RESULT_FORMAT"]
  if !ok {
    glog.Warningf("ERR: REPORT: 'REPORT_RESULT_FORMAT' don`t set\n")
    repFormat = "HTML"
  }

  repTitle, ok = (*prop)["REPORT_RESULT_TITLE"]
  if !ok {
    glog.Warningf("ERR: REPORT: 'REPORT_RESULT_TITLE' don`t set\n")
  }

  tmpl := getTemplate(repTemplate, repLang)
  if tmpl == nil {
    glog.Errorf("ERR: REPORT: '%s' don`t found", repTemplate)
    return result, false
  }

  var tpl bytes.Buffer
  
  if repFormat == "HTML" {
    if okLocal {
      var f *os.File
      f, err = os.Create(repLocalFile)
      if err != nil {
        glog.Errorf("ERR: REPORT(%s): FILE(%s): '%v'", repTemplate, repLocalFile, err)
        return result, false
      }
      err = tmpl.Execute(f, prop)
      f.Close()
    }
    if okBuffer {
      err = tmpl.Execute(&tpl, prop)
    }    
    if err != nil {
      glog.Errorf("ERR: REPORT(%s): FILE(%s): '%v'", repTemplate, repLocalFile, err)
      return result, false
    }
    if okBuffer {
      result[repBuffer] = tpl.String()
    }
  }
  // pandoc -o output.docx input.html
  
  if repFormat == "PDF" {
    var paperSize string
    var paperOrientation string
    var repCreator string
    paperSize, ok = (*prop)["REPORT_RESULT_FORMAT_PAPER_SIZE"]
    if !ok {
      glog.Warningf("ERR: REPORT: 'REPORT_RESULT_FORMAT_PAPER_SIZE' don`t set\n")
      paperSize = "A4"
    }
    paperOrientation, ok = (*prop)["REPORT_RESULT_FORMAT_PAPER_ORIENTATION"]
    if !ok {
      glog.Warningf("ERR: REPORT: 'REPORT_RESULT_FORMAT_PAPER_ORIENTATION' don`t set\n")
      paperOrientation = "P"
    }
    repCreator, ok = (*prop)["REPORT_RESULT_CREATOR"]
    if !ok {
      glog.Warningf("ERR: REPORT: 'REPORT_RESULT_CREATOR' don`t set\n")
    }
    err = tmpl.Execute(&tpl, prop)
    if err != nil {
      glog.Errorf("ERR: REPORT(%s): FILE(%s): '%v'", repTemplate, repLocalFile, err)
      return result, false
    }
    
    pdf := gofpdf.New(paperOrientation, "mm", paperSize, "")
    
    pdf.AddUTF8Font("PTAstraSans", "BI", "./fonts/PTAstraSans-BoldItalic.ttf")
    pdf.AddUTF8Font("PTAstraSans", "I", "./fonts/PTAstraSans-Italic.ttf")
    pdf.AddUTF8Font("PTAstraSans", "B", "./fonts/PTAstraSans-Bold.ttf")
    pdf.AddUTF8Font("PTAstraSans", "", "./fonts/PTAstraSans-Regular.ttf")

    pdf.AddUTF8Font("PTAstraSerif", "BI", "./fonts/PTAstraSerif-BoldItalic.ttf")
    pdf.AddUTF8Font("PTAstraSerif", "I", "./fonts/PTAstraSerif-Italic.ttf")
    pdf.AddUTF8Font("PTAstraSerif", "B", "./fonts/PTAstraSerif-Bold.ttf")
    pdf.AddUTF8Font("PTAstraSerif", "", "./fonts/PTAstraSerif-Regular.ttf")
    
    pdf.SetTitle(repTitle, true)
    pdf.SetCreator(repCreator, true)
    pdf.SetFont("PTAstraSans", "", 20)
    pdf.SetFontSize(14)
    pdf.AddPage()
    _, lineHt := pdf.GetFontSize()
    html := pdf.HTMLBasicNew()
    html.Write(lineHt, tpl.String())
    err = pdf.OutputFileAndClose(repLocalFile)
    if err != nil {
      glog.Errorf("ERR: REPORT(%s): FILE(%s): OutputFileAndClose err='%v'", repTemplate, repLocalFile, err)
      return result, false
    }
  }

  if glog.V(2) {
    glog.Infof("LOG: REPORT(%s): DONE\n", repTemplate)
  }
  return result, true
}


func getTemplate(name string, lang string) *template.Template {
  index := name + "." + lang
  i, ok := memTemplate[index]
  if ok {
    return i
  }
  var err error
  funcMap := template.FuncMap{
    "MONEY_RU": func(i string) string {
        return moneyRu(i, false)
    },
  }
  contents, err := ioutil.ReadFile(fmt.Sprintf("./templates/%s.tpl", index))
  if err != nil {
    glog.Errorf("ERR: Get Template(%s): %v", index, err)
    return nil
  }
  t := template.New(index).Funcs(funcMap)
  t, err = t.Parse(string(contents))
  if err != nil {
    glog.Errorf("ERR: Parse Template(%s): %v", index, err)
    return nil
  }
  memTemplate[index] = t
  return t
}
