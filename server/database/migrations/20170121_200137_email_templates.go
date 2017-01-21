package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type EmailTemplates_20170121_200137 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &EmailTemplates_20170121_200137{}
	m.Created = "20170121_200137"
	migration.Register("EmailTemplates_20170121_200137", m)
}

// Run the migrations
func (m *EmailTemplates_20170121_200137) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'body', '', '[body:content]', 'active', 'item', 'message', '2014-06-21 21:56:07', '2015-04-14 15:24:39' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'callout', '', '<table class="row callout">
    <tr>
        <td class="wrapper last">

            <table class="twelve columns">
                <tr>
                    <td class="panel">
                        [callout:content]
                    </td>
                    <td class="expander"></td>
                </tr>
            </table>

        </td>
    </tr>
</table>', 'active', 'item', 'main', '2014-06-15 17:45:31', '2015-04-12 01:40:03' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'contacts', '', '<table class="six columns">
    <tr>
        <td class="last right-text-pad">
        <h3>Контакты:</h3>
            [contacts:content]
        </td>
        <td class="expander"></td>
    </tr>
</table>', 'active', 'item', 'footer', '2014-06-15 17:45:20', '2015-04-12 01:40:03' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'facebook', '', '<table class="tiny-button facebook">
    <tr>
        <td>
            <a href="#">Facebook</a>
        </td>
    </tr>
</table>

<br>', 'active', 'item', 'social', '2014-06-15 17:46:59', '2015-04-13 01:50:01' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'footer', '', '<table class="row footer">
    <tr>

        <td class="wrapper ">

            [social:block]

        </td>
        <td class="wrapper last">

            [contacts:block]

        </td>

    </tr>
</table>', 'active', 'item', 'main', '2014-06-15 17:45:07', '2015-04-12 01:40:03' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'google', '', '<table class="tiny-button google-plus">
    <tr>
        <td>
            <a href="#">Google +</a>
        </td>
    </tr>
</table>

<br>', 'active', 'item', 'social', '2014-06-15 17:47:17', '2015-04-12 01:40:01' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'header', '', '<table class="row header">
    <tr>
        <td class="center" align="center">
            <center>

                <table class="container">
                    <tr>
                        <td class="wrapper last">

                            <table class="twelve columns">
                                <tr>
                                    <td class="six sub-columns">
                                        <img alt="Облачная типография, Календари-домики, Листовки и Флаеры, Карманные календари, Визитки, Буклеты, Пластиковые карты" src="data:image/png;base64,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">
                                    </td>
                                    <td class="six sub-columns last" style="text-align:right; vertical-align:middle;">
                                    </td>
                                    <td class="expander"></td>
                                </tr>
                            </table>

                        </td>
                    </tr>
                </table>

            </center>
        </td>
    </tr>
</table>', 'active', 'item', 'main', '2014-06-15 17:44:55', '2015-04-16 23:43:22' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'main', 'Основной слой', '<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="viewport" content="width=device-width"/>
  <style type="text/css">

#outlook a {
  padding:0;
}

body{
  width:100% !important;
  min-width: 100%;
  -webkit-text-size-adjust:100%;
  -ms-text-size-adjust:100%;
  margin:0;
  padding:0;
}

.ExternalClass {
  width:100%;
}

.ExternalClass,
.ExternalClass p,
.ExternalClass span,
.ExternalClass font,
.ExternalClass td,
.ExternalClass div {
  line-height: 100%;
}

#backgroundTable {
  margin:0;
  padding:0;
  width:100% !important;
  line-height: 100% !important;
}

img {
  outline:none;
  text-decoration:none;
  -ms-interpolation-mode: bicubic;
  width: auto;
  max-width: 100%;
  float: left;
  clear: both;
  display: block;
}

center {
  width: 100%;
  min-width: 580px;
}

a img {
  border: none;
}

p {
  margin: 0 0 0 10px;
}

table {
  border-spacing: 0;
  border-collapse: collapse;
}

td {
  word-break: break-word;
  -webkit-hyphens: auto;
  -moz-hyphens: auto;
  hyphens: auto;
  border-collapse: collapse !important;
}

table, tr, td {
  padding: 0;
  vertical-align: top;
  text-align: left;
}

hr {
  color: #d9d9d9;
  background-color: #d9d9d9;
  height: 1px;
  border: none;
}

/* Responsive Grid */

table.body {
  height: 100%;
  width: 100%;
}

table.container {
  width: 580px;
  margin: 0 auto;
  text-align: inherit;
}

table.row {
  padding: 0;
  width: 100%;
  position: relative;
}

table.container table.row {
  display: block;
}

td.wrapper {
  padding: 10px 20px 0 0;
  position: relative;
}

table.columns,
table.column {
  margin: 0 auto;
}

table.columns td,
table.column td {
  padding: 0 0 10px;
}

table.columns td.sub-columns,
table.column td.sub-columns,
table.columns td.sub-column,
table.column td.sub-column {
  padding-right: 10px;
}

td.sub-column, td.sub-columns {
  min-width: 0;
}

table.row td.last,
table.container td.last {
  padding-right: 0;
}

table.one { width: 30px; }
table.two { width: 80px; }
table.three { width: 130px; }
table.four { width: 180px; }
table.five { width: 230px; }
table.six { width: 280px; }
table.seven { width: 330px; }
table.eight { width: 380px; }
table.nine { width: 430px; }
table.ten { width: 480px; }
table.eleven { width: 530px; }
table.twelve { width: 580px; }

table.one center { min-width: 30px; }
table.two center { min-width: 80px; }
table.three center { min-width: 130px; }
table.four center { min-width: 180px; }
table.five center { min-width: 230px; }
table.six center { min-width: 280px; }
table.seven center { min-width: 330px; }
table.eight center { min-width: 380px; }
table.nine center { min-width: 430px; }
table.ten center { min-width: 480px; }
table.eleven center { min-width: 530px; }
table.twelve center { min-width: 580px; }

table.one .panel center { min-width: 10px; }
table.two .panel center { min-width: 60px; }
table.three .panel center { min-width: 110px; }
table.four .panel center { min-width: 160px; }
table.five .panel center { min-width: 210px; }
table.six .panel center { min-width: 260px; }
table.seven .panel center { min-width: 310px; }
table.eight .panel center { min-width: 360px; }
table.nine .panel center { min-width: 410px; }
table.ten .panel center { min-width: 460px; }
table.eleven .panel center { min-width: 510px; }
table.twelve .panel center { min-width: 560px; }

.body .columns td.one,
.body .column td.one { width: 8.333333%; }
.body .columns td.two,
.body .column td.two { width: 16.666666%; }
.body .columns td.three,
.body .column td.three { width: 25%; }
.body .columns td.four,
.body .column td.four { width: 33.333333%; }
.body .columns td.five,
.body .column td.five { width: 41.666666%; }
.body .columns td.six,
.body .column td.six { width: 50%; }
.body .columns td.seven,
.body .column td.seven { width: 58.333333%; }
.body .columns td.eight,
.body .column td.eight { width: 66.666666%; }
.body .columns td.nine,
.body .column td.nine { width: 75%; }
.body .columns td.ten,
.body .column td.ten { width: 83.333333%; }
.body .columns td.eleven,
.body .column td.eleven { width: 91.666666%; }
.body .columns td.twelve,
.body .column td.twelve { width: 100%; }

td.offset-by-one { padding-left: 50px; }
td.offset-by-two { padding-left: 100px; }
td.offset-by-three { padding-left: 150px; }
td.offset-by-four { padding-left: 200px; }
td.offset-by-five { padding-left: 250px; }
td.offset-by-six { padding-left: 300px; }
td.offset-by-seven { padding-left: 350px; }
td.offset-by-eight { padding-left: 400px; }
td.offset-by-nine { padding-left: 450px; }
td.offset-by-ten { padding-left: 500px; }
td.offset-by-eleven { padding-left: 550px; }

td.expander {
  visibility: hidden;
  width: 0;
  padding: 0 !important;
}

table.columns .text-pad,
table.column .text-pad {
  padding-left: 10px;
  padding-right: 10px;
}

table.columns .left-text-pad,
table.columns .text-pad-left,
table.column .left-text-pad,
table.column .text-pad-left {
  padding-left: 10px;
}

table.columns .right-text-pad,
table.columns .text-pad-right,
table.column .right-text-pad,
table.column .text-pad-right {
  padding-right: 10px;
}

/* Block Grid */

.block-grid {
  width: 100%;
  max-width: 580px;
}

.block-grid td {
  display: inline-block;
  padding:10px;
}

.two-up td {
  width:270px;
}

.three-up td {
  width:173px;
}

.four-up td {
  width:125px;
}

.five-up td {
  width:96px;
}

.six-up td {
  width:76px;
}

.seven-up td {
  width:62px;
}

.eight-up td {
  width:52px;
}

/* Alignment & Visibility Classes */

table.center, td.center {
  text-align: center;
}

h1.center,
h2.center,
h3.center,
h4.center,
h5.center,
h6.center {
  text-align: center;
}

span.center {
  display: block;
  width: 100%;
  text-align: center;
}

img.center {
  margin: 0 auto;
  float: none;
}

.show-for-small,
.hide-for-desktop {
  display: none;
}

/* Typography */

body, table.body, h1, h2, h3, h4, h5, h6, p, td {
  color: #222222;
  font-family: "Helvetica", "Arial", sans-serif;
  font-weight: normal;
  padding:0;
  margin: 0;
  text-align: left;
  line-height: 1.3;
}

h1, h2, h3, h4, h5, h6 {
  word-break: normal;
}

/*h1 {font-size: 40px;}*/
h1 {font-size: 30px;}
/*h2 {font-size: 36px;}*/
h2 {font-size: 26px;}
h3 {font-size: 32px;}
h4 {font-size: 28px;}
h5 {font-size: 27px;}
h6 {font-size: 20px;}
body, table.body, p, td {font-size: 14px;line-height:19px;}

p.lead, p.lede, p.leed {
  font-size: 18px;
  line-height:21px;
}

p {
  margin-bottom: 10px;
}

small {
  font-size: 10px;
}

a {
  color: #2ba6cb;
  text-decoration: none;
}

a:hover {
  color: #2795b6 !important;
}

a:active {
  color: #2795b6 !important;
}

a:visited {
  color: #2ba6cb !important;
}

h1 a,
h2 a,
h3 a,
h4 a,
h5 a,
h6 a {
  color: #2ba6cb;
}

h1 a:active,
h2 a:active,
h3 a:active,
h4 a:active,
h5 a:active,
h6 a:active {
  color: #2ba6cb !important;
}

h1 a:visited,
h2 a:visited,
h3 a:visited,
h4 a:visited,
h5 a:visited,
h6 a:visited {
  color: #2ba6cb !important;
}

/* Panels */

.panel {
  background: #f2f2f2;
  border: 1px solid #d9d9d9;
  padding: 10px !important;
}

.sub-grid table {
  width: 100%;
}

.sub-grid td.sub-columns {
  padding-bottom: 0;
}

/* Buttons */

table.button,
table.tiny-button,
table.small-button,
table.medium-button,
table.large-button {
  width: 100%;
  overflow: hidden;
}

table.button td,
table.tiny-button td,
table.small-button td,
table.medium-button td,
table.large-button td {
  display: block;
  width: auto !important;
  text-align: center;
  background: #2ba6cb;
  border: 1px solid #2284a1;
  color: #ffffff;
  padding: 8px 0;
}

table.tiny-button td {
  padding: 5px 0 4px;
}

table.small-button td {
  padding: 8px 0 7px;
}

table.medium-button td {
  padding: 12px 0 10px;
}

table.large-button td {
  padding: 21px 0 18px;
}

table.button td a,
table.tiny-button td a,
table.small-button td a,
table.medium-button td a,
table.large-button td a {
  font-weight: bold;
  text-decoration: none;
  font-family: Helvetica, Arial, sans-serif;
  color: #ffffff;
  font-size: 16px;
}

table.tiny-button td a {
  font-size: 12px;
  font-weight: normal;
}

table.small-button td a {
  font-size: 16px;
}

table.medium-button td a {
  font-size: 20px;
}

table.large-button td a {
  font-size: 24px;
}

table.button:hover td,
table.button:visited td,
table.button:active td {
  background: #2795b6 !important;
}

table.button:hover td a,
table.button:visited td a,
table.button:active td a {
  color: #fff !important;
}

table.button:hover td,
table.tiny-button:hover td,
table.small-button:hover td,
table.medium-button:hover td,
table.large-button:hover td {
  background: #2795b6 !important;
}

table.button:hover td a,
table.button:active td a,
table.button td a:visited,
table.tiny-button:hover td a,
table.tiny-button:active td a,
table.tiny-button td a:visited,
table.small-button:hover td a,
table.small-button:active td a,
table.small-button td a:visited,
table.medium-button:hover td a,
table.medium-button:active td a,
table.medium-button td a:visited,
table.large-button:hover td a,
table.large-button:active td a,
table.large-button td a:visited {
  color: #ffffff !important;
}

table.secondary td {
  background: #e9e9e9;
  border-color: #d0d0d0;
  color: #555;
}

table.secondary td a {
  color: #555;
}

table.secondary:hover td {
  background: #d0d0d0 !important;
  color: #555;
}

table.secondary:hover td a,
table.secondary td a:visited,
table.secondary:active td a {
  color: #555 !important;
}

table.success td {
  background: #5da423;
  border-color: #457a1a;
}

table.success:hover td {
  background: #457a1a !important;
}

table.alert td {
  background: #c60f13;
  border-color: #970b0e;
}

table.alert:hover td {
  background: #970b0e !important;
}

table.radius td {
  -webkit-border-radius: 3px;
  -moz-border-radius: 3px;
  border-radius: 3px;
}

table.round td {
  -webkit-border-radius: 500px;
  -moz-border-radius: 500px;
  border-radius: 500px;
}

/* Outlook First */

body.outlook p {
  display: inline !important;
}

/*  Media Queries */

@media only screen and (max-width: 600px) {

  table[class="body"] img {
    width: auto !important;
    height: auto !important;
  }

  table[class="body"] center {
    min-width: 0 !important;
  }

  table[class="body"] .container {
    width: 95% !important;
  }

  table[class="body"] .row {
    width: 100% !important;
    display: block !important;
  }

  table[class="body"] .wrapper {
    display: block !important;
    padding-right: 0 !important;
  }

  table[class="body"] .columns,
  table[class="body"] .column {
    table-layout: fixed !important;
    float: none !important;
    width: 100% !important;
    padding-right: 0 !important;
    padding-left: 0 !important;
    display: block !important;
  }

  table[class="body"] .wrapper.first .columns,
  table[class="body"] .wrapper.first .column {
    display: table !important;
  }

  table[class="body"] table.columns td,
  table[class="body"] table.column td {
    width: 100% !important;
  }

  table[class="body"] .columns td.one,
  table[class="body"] .column td.one { width: 8.333333% !important; }
  table[class="body"] .columns td.two,
  table[class="body"] .column td.two { width: 16.666666% !important; }
  table[class="body"] .columns td.three,
  table[class="body"] .column td.three { width: 25% !important; }
  table[class="body"] .columns td.four,
  table[class="body"] .column td.four { width: 33.333333% !important; }
  table[class="body"] .columns td.five,
  table[class="body"] .column td.five { width: 41.666666% !important; }
  table[class="body"] .columns td.six,
  table[class="body"] .column td.six { width: 50% !important; }
  table[class="body"] .columns td.seven,
  table[class="body"] .column td.seven { width: 58.333333% !important; }
  table[class="body"] .columns td.eight,
  table[class="body"] .column td.eight { width: 66.666666% !important; }
  table[class="body"] .columns td.nine,
  table[class="body"] .column td.nine { width: 75% !important; }
  table[class="body"] .columns td.ten,
  table[class="body"] .column td.ten { width: 83.333333% !important; }
  table[class="body"] .columns td.eleven,
  table[class="body"] .column td.eleven { width: 91.666666% !important; }
  table[class="body"] .columns td.twelve,
  table[class="body"] .column td.twelve { width: 100% !important; }

  table[class="body"] td.offset-by-one,
  table[class="body"] td.offset-by-two,
  table[class="body"] td.offset-by-three,
  table[class="body"] td.offset-by-four,
  table[class="body"] td.offset-by-five,
  table[class="body"] td.offset-by-six,
  table[class="body"] td.offset-by-seven,
  table[class="body"] td.offset-by-eight,
  table[class="body"] td.offset-by-nine,
  table[class="body"] td.offset-by-ten,
  table[class="body"] td.offset-by-eleven {
    padding-left: 0 !important;
  }

  table[class="body"] table.columns td.expander {
    width: 1px !important;
  }

  table[class="body"] .right-text-pad,
  table[class="body"] .text-pad-right {
    padding-left: 10px !important;
  }

  table[class="body"] .left-text-pad,
  table[class="body"] .text-pad-left {
    padding-right: 10px !important;
  }

  table[class="body"] .hide-for-small,
  table[class="body"] .show-for-desktop {
    display: none !important;
  }

  table[class="body"] .show-for-small,
  table[class="body"] .hide-for-desktop {
    display: inherit !important;
  }
}

  </style>
  <style type="text/css">

    table.facebook td {
      background: #3b5998;
      border-color: #2d4473;
    }

    table.facebook:hover td {
      background: #2d4473 !important;
    }

    table.twitter td {
      background: #00acee;
      border-color: #0087bb;
    }

    table.twitter:hover td {
      background: #0087bb !important;
    }

    table.google-plus td {
      background-color: #DB4A39;
      border-color: #CC0000;
    }

    table.google-plus:hover td {
      background: #CC0000 !important;
    }

    .template-label {
      color: #ffffff;
      font-weight: bold;
      font-size: 11px;
    }

    .callout .wrapper {
      padding-bottom: 20px;
    }

    .callout .panel {
      background: #ECF8FF;
      border-color: #b9e5ff;
    }

    .header {
      background: #394e63;
      min-height:100px;
    }

    .footer .wrapper {
      background: #ebebeb;
    }

    .footer h5 {
      padding-bottom: 10px;
    }

    table.columns .text-pad {
      padding-left: 10px;
      padding-right: 10px;
    }

    table.columns .left-text-pad {
      padding-left: 10px;
    }

    table.columns .right-text-pad {
      padding-right: 10px;
    }

    @media only screen and (max-width: 600px) {

      table[class="body"] .right-text-pad {
        padding-left: 10px !important;
      }

      table[class="body"] .left-text-pad {
        padding-right: 10px !important;
      }
    }

  </style>
</head>
<body>
  <table class="body">
    <tr>
      <td class="center" align="center" valign="top">
        <center>

            [header:block]

          <table class="container">
            <tr>
              <td>

                [message:block]

                [callout:block]

                [footer:block]

                [privacy:block]

              <!-- container end below -->
              </td>
            </tr>
          </table>

        </center>
      </td>
    </tr>
  </table>
</body>
</html>', 'active', 'item', NULL, '2014-06-15 17:44:12', '2015-04-16 23:47:23' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'message', '', '<table class="row">
  <tr>
      <td class="wrapper last">

          <table class="twelve columns">
              <tr>
                  <td>
                      [title:block]
                      [body:block]
                  </td>
                  <td class="expander"></td>
              </tr>
          </table>

      </td>
  </tr>
</table>', 'active', 'item', 'main', '2014-06-20 09:15:51', '2015-04-12 01:40:00' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'password_reset', 'Восстановление пароля', '{"items":["body","header"],"title":"Смена логина или пароля для [user:name:last] [user:name:first] на сайте [site:name]","fields":[{"name":"body","value":"[user:name:last] [user:name:first],<br><br>Запрос на сброс пароля для вашего аккаунта был сделан на сайте [site:name]. <br><br>Вы можете сейчас войти на сайт, кликнув на ссылку или скопировав и вставив её в браузер: <br><br>[user:one-time-login-url] <br><br>Это одноразовая ссылка для входа и она перебросит Вас на страницу, где Вы сможете изменить пароль. Ссылка истекает через 1 сутки и ничего не случится, если она не будет использована. <br><br>-- Команда сайта [site:name]"}]}', 'active', 'template', NULL, '2014-06-20 18:35:00', '2015-04-16 17:27:39' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'privacy', '', '<table class="row">
    <tr>
        <td class="wrapper last">

            <table class="twelve columns">
                <tr>
                    <td align="center">
                        <center>
                            [privacy:content]
                        </center>
                    </td>
                    <td class="expander"></td>
                </tr>
            </table>
        </td>
    </tr>
</table>', 'active', 'item', 'main', '2014-06-15 17:45:42', '2015-04-12 01:38:28' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'register_admin_created', 'Добро пожаловать (новый пользователь создан администратором)', '{"items":["header","body"],"title":"Администратор создал для вас учётную запись","fields":[{"name":"body","value":"[user:name:last] [user:name:first],<br><br>Администратор системы [site:name] создал для вас аккаунт. Можете войти в систему, кликнув на ссылку или скопировав и вставив её в адресную строку браузера:<br><br>[user:one-time-login-url] <br><br>Эта одноразовая ссылка для входа в систему направит вас на страницу задания своего пароля.<br><br>После установки пароля вы сможете входить в систему через страницу<br>[site:login-url]<br><br>-- Команда [site:name]"}]}', 'active', 'template', '', '2014-06-20 18:27:05', '2017-01-15 17:03:11' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'social', '', '<table class="six columns">
    <tr>
        <td class="left-text-pad">

            <h3>Мы в сети:</h3>

            [facebook:block]
            [twitter:block]
            [vk:block]
            [google:block]

</td>
        <td class="expander"></td>
    </tr>
</table>', 'active', 'item', 'footer', '2014-06-15 17:46:47', '2015-06-23 14:10:14' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'status_activated', 'Учётная запись активирована', '{"items":["header","body"],"title":"Детали учётной записи для [user:name:last] [user:name:first] на [site:name] (одобрено)","fields":[{"name":"body","value":"[user:name:last] [user:name:first],<br><br>Ваш аккаунт на сайте [site:name] был активирован.<br><br>Вы можете войти на сайт, кликнув на ссылку или скопировав и вставив её в Ваш браузер: <br><br>[user:one-time-login-url] <br><br>Это одноразовая ссылка для входа и она перебросит Вас на страницу, где Вы сможете установить свой пароль.<br><br>После установки Вашего пароля, вы сможете заходить на странице [site:login-url].<br><br>-- Команда сайта [site:name]"}]}', 'active', 'template', '', '2014-06-20 18:31:06', '2017-01-15 06:07:52' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'title', '', '<h1>[title:content]</h1>
', 'active', 'item', 'message', '2014-06-21 21:54:48', '2015-04-14 15:20:18' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'twitter', '', '<table class="tiny-button twitter">
    <tr>
        <td>
            <a href="#">Twitter</a>
        </td>
    </tr>
</table>

<br>', 'active', 'item', 'social', '2014-06-15 17:47:08', '2015-04-13 01:47:52' );`)
	m.SQL(`INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( 'vk', '', '<table class="tiny-button vk">
    <tr>
        <td>
            <a href="#">vk</a>
        </td>
    </tr>
</table>

<br>', 'active', 'item', 'social', '2014-06-15 17:50:15', '2015-04-13 01:51:30' );`)
}

// Reverse the migrations
func (m *EmailTemplates_20170121_200137) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DELETE FROM `email_templates`")
}
