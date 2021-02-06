package main

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
)


/////////////////////////
// TESTS
/////////////////////////
func TestMakeReport(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
  
  tmpl := getTemplate("unknown", "ru")
  if tmpl != nil {
    t.Error(
      "For", "ERR: GET unknown template",
      "expected", "+",
      "got", tmpl,
    )
  }
  
  cfg := ConfigInfo{DefaultLang: "ru"}
  
  report_settings := ReportInfo{}
  prop := map[string]string{"REPORT_TEMPLATE": "unknown",
              "REPORT_RESULT_LOCAL_FILE": "./storage/invoice.ru.1.html",
              "ACCOUNT_TO_NAME": "ООО \"Получатель\"",
              "ACCOUNT_TO_INDEX": "127282",
              "ACCOUNT_TO_CITY": "Москва",
              "ACCOUNT_FROM_CITY": "Москва",
              "PAYMENT_SUM_WITH_VAT": "105.23",
             }
  res, ok := makeReport(&cfg, &report_settings, &prop)
  assert.Equal(t, false, ok)

  prop = map[string]string{"REPORT_TEMPLATE": "invoice",
              "ACCOUNT_TO_NAME": "ООО \"Получатель\"",
              "ACCOUNT_TO_INDEX": "127282",
              "ACCOUNT_TO_CITY": "Москва",
              "ACCOUNT_FROM_CITY": "Москва",
              "PAYMENT_SUM_WITH_VAT": "105.23",
             }
  res, ok = makeReport(&cfg, &report_settings, &prop)
  assert.Equal(t, false, ok)

  prop = map[string]string{}
  res, ok = makeReport(&cfg, &report_settings, &prop)
  assert.Equal(t, false, ok)

  prop = map[string]string{"REPORT_TEMPLATE": "invoice", "REPORT_RESULT_LOCAL_FILE": "./storage/invoice.ru.1.html",
              "ACCOUNT_TO_NAME": "ООО \"Получатель\"",
              "ACCOUNT_TO_INDEX": "127282",
              "ACCOUNT_TO_CITY": "Москва",
              "ACCOUNT_FROM_CITY": "Москва",
              "PAYMENT_SUM_WITH_VAT": "105.23",
             }
  res, ok = makeReport(&cfg, &report_settings, &prop)
  assert.Equal(t, true, ok)

  prop2 := map[string]string{"REPORT_TEMPLATE": "invoice", "REPORT_RESULT_LOCAL_FILE": "./storage/invoice.ru.1.pdf",
              "ACCOUNT_TO_NAME": "ООО \"Получатель\"",
              "ACCOUNT_TO_INDEX": "127282",
              "ACCOUNT_TO_CITY": "Москва",
              "ACCOUNT_FROM_CITY": "Москва",
              "PAYMENT_SUM_WITH_VAT": "105.23",
              "REPORT_RESULT_FORMAT": "PDF",
             }
  res, ok = makeReport(&cfg, &report_settings, &prop2)
  assert.Equal(t, true, ok)
  
  prop2["REPORT_RESULT_FORMAT"] = "HTML"
  prop2["REPORT_RESULT_BUFFER"] = "MAIL_BODY"
  res, ok = makeReport(&cfg, &report_settings, &prop2)
  assert.Equal(t, true, ok)
  
  res_need := map[string]string{"MAIL_BODY": "<!doctype html>\n<html>\n<head>\n    <title>Бланк \"Счет на оплату\"</title>\n    <meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\">\n    <style>\n        body { width: 210mm; margin-left: auto; margin-right: auto; border: 1px #efefef solid; font-size: 11pt;}\n        table.invoice_bank_rekv { cellspacing: 0; border-collapse: collapse; padding: 0; border: 0.1px solid black; }\n        table.invoice_bank_rekv > tbody > tr > td, table.invoice_bank_rekv > tr > td { border-collapse: collapse; border: 1px solid black; }\n        table.invoice_items { border: 1px solid black; border-collapse: collapse; padding: 0; cellspacing: 0; }\n        table.invoice_items td, table.invoice_items th { border-collapse: collapse; border: 1px solid black;}\n    </style>\n</head>\n<body>\n<table width=\"100%\">\n    <tr>\n        \n        <td >\n            Внимание! Оплата данного счета означает согласие с условиями поставки товара. Уведомление об оплате  обязательно, в противном случае не гарантируется наличие товара на складе. Товар отпускается по факту прихода денег на р/с Поставщика, самовывозом, при наличии доверенности и паспорта.\n        </td>\n    </tr>\n\n</table>\n  \n<br>\n<br>\n\n<table width=\"90%\" cellpadding=\"2\" cellspacing=\"0\" class=\"invoice_bank_rekv\">\n    <tr>\n        <td colspan=\"2\" rowspan=\"2\" style=\"min-height:13mm; width: 105mm;\">\n            <table width=\"100%\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" style=\"height: 13mm;\">\n                <tr>\n                    <td valign=\"top\">\n                        <div></div>\n                    </td>\n                </tr>\n                <tr>\n                    <td valign=\"bottom\" style=\"height: 3mm;\">\n                        <div style=\"font-size:10pt;\">Банк получателя        </div>\n                    </td>\n                </tr>\n            </table>\n        </td>\n        <td style=\"min-height:7mm;height:auto; width: 25mm;\">\n            <div>БИK</div>\n        </td>\n        <td rowspan=\"2\" style=\"vertical-align: top; width: 60mm;\">\n            <div></div>\n            <div></div>\n        </td>\n    </tr>\n    <tr>\n        <td style=\"width: 25mm;\">\n            <div>Сч. №</div>\n        </td>\n    </tr>\n    <tr>\n        <td style=\"min-height:6mm; height:auto; width: 50mm;\">\n            <div>ИНН </div>\n        </td>\n        <td style=\"min-height:6mm; height:auto; width: 55mm;\">\n            <div>КПП </div>\n        </td>\n        <td rowspan=\"2\" style=\"min-height:19mm; height:auto; vertical-align: top; width: 25mm;\">\n            <div>Сч. №</div>\n        </td>\n        <td rowspan=\"2\" style=\"min-height:19mm; height:auto; vertical-align: top; width: 60mm;\">\n            <div></div>\n        </td>\n    </tr>\n    <tr>\n        <td colspan=\"2\" style=\"min-height:13mm; height:auto;\">\n\n            <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" style=\"height: 13mm; width: 105mm;\">\n                <tr>\n                    <td valign=\"top\">\n                        <div></div>\n                    </td>\n                </tr>\n                <tr>\n                    <td valign=\"bottom\" style=\"height: 3mm;\">\n                        <div style=\"font-size: 10pt;\">Получатель</div>\n                    </td>\n                </tr>\n            </table>\n\n        </td>\n    </tr>\n</table>\n<br/>\n\n<div style=\"font-weight: bold; font-size: 16pt; padding-left:5px;\">\n    Счет № от </div>\n<br/>\n\n<div style=\"background-color:#000000; width:100%; font-size:1px; height:2px;\">&nbsp;</div>\n\n<table width=\"100%\">\n    <tr>\n        <td style=\"width: 30mm;\">\n            <div style=\" padding-left:2px;\">Поставщик:</div>\n        </td>\n        <td>\n            <div style=\"font-weight:bold;  padding-left:2px;\">\n                ООО &#34;Получатель&#34;, ИНН , КПП , , , , Телефон: </div>\n        </td>\n    </tr>\n    <tr>\n        <td style=\"width: 30mm;\">\n            <div style=\" padding-left:2px;\">Покупатель:</div>\n        </td>\n        <td>\n            <div style=\"font-weight:bold;  padding-left:2px;\">\n                , ИНН , КПП , , , , Телефон: </div>\n        </td>\n    </tr>\n</table>\n\n\n<table class=\"invoice_items\" width=\"100%\" cellpadding=\"2\" cellspacing=\"0\">\n    <tr>\n        <th style=\"width:13mm;text-align:center;\">№</th>\n        <th style=\"width:70mm;text-align:center;\">Товар</th>\n        <th style=\"width:20mm;text-align:center;\">Кол-во</th>\n        <th style=\"width:17mm;text-align:center;\">Ед.</th>\n        <th style=\"width:27mm;text-align:center;\">Цена</th>\n        <th style=\"width:27mm;text-align:center;\">Сумма</th>\n    </tr>\n      <tr>\n          <td style=\"text-align:center;\"></td>\n          <td></td>\n          <td style=\"text-align:center;\"></td>\n          <td style=\"text-align:center;\"></td>\n          <td style=\"text-align:right;\"></td>\n          <td style=\"text-align:right;\"></td>\n      </tr>\n</table>\n\n<table border=\"0\" width=\"100%\" cellpadding=\"1\" cellspacing=\"1\">\n    <tr>\n        <td style=\"width:100mm;\">&nbsp;</td>\n        <td style=\"width:47mm; font-weight:bold;  text-align:right;\">Итого:</td>\n        <td style=\"width:27mm; font-weight:bold;  text-align:right;\"></td>\n    </tr>\n    <tr>\n        <td style=\"width:100mm;\">&nbsp;</td>\n        <td style=\"width:47mm; font-weight:bold;  text-align:right;\">В том числе НДС:</td>\n        <td style=\"width:27mm; font-weight:bold;  text-align:right;\">105.23</td>\n    </tr>\n</table>\n\n<br />\n<div>\nВсего наименований 0 на сумму 105.23 рублей.<br />\n<b>сто пять рублей 23 копейки</b></div>\n<br /><br />\n<div style=\"background-color:#000000; width:100%; font-size:1px; height:2px;\">&nbsp;</div>\n<br/>\n\n<div> ______________________ ()</div>\n<br/>\n\n<div>Главный бухгалтер ______________________ ()</div>\n<br/>\n\n<div style=\"width: 85mm;text-align:center;\">М.П.</div>\n<br/>\n\n\n<div style=\"width:800px;text-align:left;font-size:10pt;\">Счет действителен к оплате в течении пяти дней.</div>\n\n</body>\n</html>\n",
             }
  assert.Equal(t, res_need, res)

}
