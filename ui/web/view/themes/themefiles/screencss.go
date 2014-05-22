// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themefiles

const ScreenCss = `
html {
    min-height: 100%;
    font-size: 100%;
    overflow-y: scroll;
    -webkit-text-size-adjust: 100%;
    -ms-text-size-adjust: 100%;
}

body {
    color: #444;
    font-family: Georgia, Palatino, 'Palatino Linotype', Times, 'Times New Roman', "Hiragino Sans GB", "STXihei", "微软雅黑", serif;
    font-size: 12px;
    line-height: 1.5em;
    background: #fefefe;
    width: 75%;
    min-height: 100%;
    margin: 0px auto 1em auto;
    padding: 0 1em;
    outline: 1300px solid #FAFAFA;
}

.cleaner {
    clear: both;
}

a {
    color: #0645ad;
    text-decoration: none;
}

a:visited {
    color: #0b0080;
}

a:hover {
    color: #06e;
}

a:active {
    color: #faa700;
}

a:focus {
    outline: thin dotted;
}

a:hover, a:active {
    outline: 0;
}

span.backtick {
    border: 1px solid #EAEAEA;
    border-radius: 3px;
    background: #F8F8F8;
    padding: 0 3px 0 3px;
}

::-moz-selection {
    background: rgba(255,255,0,0.3);
    color: #000;
}

::selection {
    background: rgba(255,255,0,0.3);
    color: #000;
}

a::-moz-selection {
    background: rgba(255,255,0,0.3);
    color: #0645ad;
}

a::selection {
    background: rgba(255,255,0,0.3);
    color: #0645ad;
}

p {
    margin: 1em 0;
}

img {
    max-width: 100%;
}

h1,h2,h3,h4,h5,h6 {
    font-weight: normal;
    color: #111;
    line-height: 1em;
}

h4,h5,h6 {
    font-weight: bold;
}

h1 {
    font-size: 2.5em;
    margin: 0;
}

h2 {
    font-size: 2em;
    border-bottom: 1px solid silver;
    padding-bottom: 5px;
}

h3 {
    font-size: 1.5em;
}

h4 {
    font-size: 1.2em;
}

h5 {
    font-size: 1em;
}

h6 {
    font-size: 0.9em;
}

blockquote {
    color: #666666;
    margin: 0;
    padding-left: 3em;
    border-left: 0.5em #EEE solid;
}

hr {
    display: block;
    height: 2px;
    border: 0;
    border-top: 1px solid #aaa;
    border-bottom: 1px solid #eee;
    margin: 1em 0;
    padding: 0;
}

pre , code, kbd, samp {
    color: #000;
    font-family: monospace;
    font-size: 0.88em;
    border-radius: 3px;
    background-color: #F8F8F8;
    border: 1px solid #CCC;
}

pre {
    white-space: pre;
    white-space: pre-wrap;
    word-wrap: break-word;
    padding: 5px 12px;
}

pre code {
    border: 0px !important;
    padding: 0;
}

code {
    padding: 0 3px 0 3px;
}

b, strong {
    font-weight: bold;
}

dfn {
    font-style: italic;
}

ins {
    background: #ff9;
    color: #000;
    text-decoration: none;
}

mark {
    background: #ff0;
    color: #000;
    font-style: italic;
    font-weight: bold;
}

sub, sup {
    font-size: 75%;
    line-height: 0;
    position: relative;
    vertical-align: baseline;
}

sup {
    top: -0.5em;
}

sub {
    bottom: -0.25em;
}

ul, ol {
    margin: 0;
    padding: 0 0 0 2em;
}

li p:last-child {
    margin: 0;
}

ul.tree, ul.tree ul {
    list-style-type: none;
    background: url(tree-vertical-line.png) repeat-y;
    margin: 0;
    padding: 0;
}

ul.tree ul {
    margin-left: 10px;
}

ul.tree li {
    margin: 0;
    padding: 0 12px;
    line-height: 20px;
    background: url(tree-node.png) no-repeat;
}

ul.tree li:last-child {
    background: #fff url(tree-last-node.png) no-repeat;
}

dd {
    margin: 0 0 0 2em;
}

img {
    border: 0;
    -ms-interpolation-mode: bicubic;
    vertical-align: middle;
}

table {
    border-collapse: collapse;
    border-spacing: 0;
}

td {
    vertical-align: top;
}

body>footer {
    margin: 2em 0 0 0;
    border-top: 1px solid #eee;
}

body>footer>nav {
    text-align: right;
}

body>footer>nav>ul {
    list-style: none;
    margin: 0;
    padding: 0;
}

body>footer>nav>ul>li {
    display: inline;
    margin-left: 1em;
}

body>footer>nav>ul>li>a {
    color: #000000;
    font-size: 0.8em;
}

body>nav.toplevel {
    text-align: right;
    font-size: 1.1em;
    margin: 0 0 1.2em 0;
    padding: 0;
    clear: both;
}

body>nav.search {
    float: right;
}

body>nav.search .typeahead {
    width: 20em;
    font-family: "Helvetia", "Verdana", "Sans-Serif";
    color: #000000;
    text-align: left;
    padding: 0px 10px;
}

.typeahead,
.tt-query,
.tt-hint {
  border: 2px solid #ccc;
  -webkit-border-radius: 8px;
     -moz-border-radius: 8px;
          border-radius: 8px;
  outline: none;
}

.typeahead {
  background-color: #fff;
}

.typeahead:focus {
  border: 2px solid #0097cf;
}

.tt-query {
  -webkit-box-shadow: inset 0 1px 1px rgba(0, 0, 0, 0.075);
     -moz-box-shadow: inset 0 1px 1px rgba(0, 0, 0, 0.075);
          box-shadow: inset 0 1px 1px rgba(0, 0, 0, 0.075);
}

.tt-hint {
  color: #999
}

.tt-dropdown-menu {
  width: 422px;
  margin-top: 12px;
  padding: 8px 0;
  background-color: #fff;
  border: 1px solid #ccc;
  border: 1px solid rgba(0, 0, 0, 0.2);
  -webkit-border-radius: 8px;
     -moz-border-radius: 8px;
          border-radius: 8px;
  -webkit-box-shadow: 0 5px 10px rgba(0,0,0,.2);
     -moz-box-shadow: 0 5px 10px rgba(0,0,0,.2);
          box-shadow: 0 5px 10px rgba(0,0,0,.2);
}

.tt-suggestion {
  padding: 3px 20px;
  font-size: 18px;
  line-height: 24px;
}

.tt-suggestion.tt-cursor {
  color: #fff;
  background-color: #0097cf;

}

.tt-suggestion p {
  margin: 0;
}

body>nav.toplevel>ul {
    display: inline;
    list-style: none;
    margin: 0;
    padding: 0;
}

body>nav.toplevel>ul>li {
    display: inline-block;
    border: 1px solid #000000;
    background-color: #FFFFFF;
    margin: -1px -5px 0 -2px;
    padding: 0px 10px;
    white-space: nowrap;
}

body>nav.toplevel>ul>li:last-child {
    border-right: 1px solid #000000;
    clear: both;
}

body>nav.toplevel>ul>li>a {
    font-family: "Helvetia", "Verdana", "Sans-Serif";
    color: #000000;
}

body>nav.toplevel>ul>li>a:hover {
    color: #06e;
}

body>nav.breadcrumb {
    clear: both;
    font-size: 0.8em;
}

body>nav.breadcrumb>ul {
    list-style: none;
    margin: 0;
    padding: 0;
}

body>nav.breadcrumb>ul>li {
    display: inline;
}

body>nav.breadcrumb>ul>li:after {
    content: "»";
}

body>nav.breadcrumb>ul>li:last-child:after {
    content: "";
}

body>nav.breadcrumb>ul>li>a {
    font-family: "Helvetia", "Verdana", "Sans-Serif";
}

article>.description {
    font-size: 1.2em;
    margin: 0 0 1.5em 0;
}

article .attachments {
    margin: 1.5em 0 0 0;
    font-size: 1.0em;
}

article .attachments>header {
    font-weight: bold;
}

article .attachments ul.tree li {
    color: #369;
    font-weight: bold;
}

article .attachments ul.tree li a {
    color: #369;
    font-weight: bold;
}

article>.tags {
   margin: 1.5em 0 0 0;    
}

article>.tags>header {
    float: left;
    font-size: 1.2em;
}

article>.tags>.tags {
    list-style-type: none;

}

article>.tags>.tags>.tag {
    display: inline-block;
    margin: 0 0 0 0.5em;
}

article>.tags>.tags>.tag>a {
    color: #FFF;
    background-color: #000;
    line-height: 1.0em;
    font-size: 1.0em;
    padding: 3px 6px;
}

article>.childs {
    margin: 2em 0 0 0;
}

article>.childs>.list {
    list-style: none;
    padding: 5px 0 5px 0;
    margin: 0 0 0 1.2em;    
}

article>.childs>.list>.child {
    margin: 0 0 1.2em 0;
}

article>.childs>.list>.child:nth-child(odd) {
    background-color:#eee;
}

article>.childs>.list>.child:nth-child(even) {
    background-color:transparent;
}

.imagegallery {
    margin: 2em 0;
}

.imagegallery>header {
    font-size: 1.2em;
}

.imagegallery ol {
    list-style: none;
    margin-left: 0;
}

.filelinks>header {
    font-size: 1.2em;
}

.filelinks ul {
    margin-left: 0;
}

.csv {
    margin: 2em 0 0 2em;
    overflow: auto;
}

.csv>header {
    font-size: 1.2em;
}

.csv>table
{
    font-family: "Lucida Sans Unicode", "Lucida Grande", Sans-Serif;
    font-size: 1.0em;
    margin: 10px 45px 10px 45px;
    text-align: left;
    border-collapse: collapse;
    border: 1px solid #69c;
}

.csv>table thead
{
    padding: 12px 17px 12px 17px;
    font-weight: normal;
    font-size: 1.2em;
    color: #039;
    border-bottom: 1px dashed #69c;
}

.csv>table td
{
    padding: 7px 17px 7px 17px;
    color: #669;
}

.csv>table tbody tr:hover td
{
    color: #339;
    background: #d0dafd;
}

.pdf {
    width: 80%;
    max-height: 80%;
    margin: 2em 0;
}

.pdf>header {
    font-size: 1.2em;
}

.pdf .metadata {
    float: left;
    font-size: 1.0em;
    margin: 0 15px 0 0;
}

.pdf .metadata .entry {
    margin: 0 0 10px 0;
}

.pdf .previewarea {
    display: block;
    overflow: auto;
    box-shadow: 0 0 10px #000000;
}

.pdf .previewarea canvas {
    float: none;
}

section.video {
    margin: 2em 0;
}

section.video>header {
    font-size: 1.2em;
}

.audio>header {
    font-size: 1.2em;
}

.presentation nav {
    float: left;
    margin: 2em 0 1.2em 0;
    width: 100%;
    text-align: center;
}

.presentation nav .controls {
    float: left;
}

.presentation nav .pager {
    display: inline;
    cursor: default;
}

.presentation nav .jumper {
    float: right;
}

.presentation .content {
    clear: left;
}

.presentation .slide {
    float: none;
    box-shadow: 0 0 10px #000000;
    padding: 10px;
}

.filepreview {
    margin: 2em 0;
}

.filepreview>header {
    font-size: 1.2em;
}

/* first level */
.sitemap>.content>ol {
    list-style: none;
    padding: 0;
    margin: 15px 0 0 0;
}

/* second level */
.sitemap>.content>ol>li>ol {
    list-style-type: lower-alpha;
    padding: 0 40px;
}

.sitemap>.content>ol>li>ol>li {
    padding: 15px 0;
}

/* third, forth, n-th level */
.sitemap>.content>ol>li>ol>li ol {
    list-style-type: decimal;
}

.tagcloud {
    margin: 0px 20px 35px 20px;
}

.tagcloud>.tags>span {
    line-height: 1.0em;
    padding: 3px 6px;
    margin: 0px 5px;
}

.tagmap>.content>.tagcloud>.tags>span>a {
    color: #00;
}

.tagmap>.content>.tagcloud>.tags>span.level-1 {
    font-size: 2em;
}

.tagmap>.content>.tagcloud>.tags>span.level-2 {
    font-size: 1.8em;
}

.tagmap>.content>.tagcloud>.tags>span.level-3 {
    font-size: 1.6em;
}

.tagmap>.content>.tagcloud>.tags>span.level-4 {
    font-size: 1.4em;
}

.tagmap>.content>.tagcloud>.tags>span.level-5 {
    font-size: 1.2em;
}

.tagmap>.content>.tagcloud>.tags>span.level-6 {
    font-size: 1.0em;
}

.tagmap>.content>.tags {
    list-style-type: none;
}

.tagmap>.content>.tags>.tag {
    margin: 0 0 1.5em 0;
}

.tagmap>.content>.tags>.tag>a {
    color: #FFF;
    background-color: #000;
    line-height: 1.2em;
    font-size: 1.2em;
    padding: 3px 6px;
}

.tagmap>.content>.tags>.tag>.childs {
    list-style-type: none;
}

.tagmap>.content>.tags>.tag>.childs>.child {
    margin: 0.7em 0 0;
}

.search>.content>header {
    margin: 10px 0 10px 0;
}

.search>.content>ol {
}

.search>.content>ol>li {
    margin: 0 0 15px 0;
    width: 35em;
}

.search>.content>ol>li>.title {
    margin: 0px;
    font-size: 1.1em;
}

.search>.content>ol>li>.description {
    margin: 0px;
    font-size: 1.0em;
}

.search>.content>ol>li>.path {
    margin: 0px;
    font-size: 0.8em;
}

@media only screen and (max-height: 500px) {
    body>article {
        min-height: 300px;
    }

    .presentation .slide {
        min-height: 220px;
    }
}

@media only screen and (min-height: 500px) {
    body>article {
        min-height: 500px;
    }

    .presentation .slide {
        min-height: 320px;
    }
}

@media only screen and (min-height: 600px) {
    body>article {
        min-height: 600px;
    }

    .presentation .slide {
        min-height: 420px;
    }
}

@media only screen and (min-height: 768px) {
    body>article {
        min-height: 768px;
    }

    .presentation .slide {
        min-height: 520px;
    }
}

@media only screen and (max-width: 480px) {
    body {
        font-size: 12px;
        width: 95%;
    }

    .presentation nav {
        margin: 25px 0 15px 0;
        width: 100%;
    }    

    .presentation nav .pager {
        float: right;
    }    

    .presentation nav .jumper {
        display: none;
    }    
}

@media only screen and (min-width: 480px) {
    body {
        font-size: 14px;
        width: 95%;
    }
}

@media only screen and (min-width: 768px) {
    body {
        font-size: 16px;
        width: 95%;
    }
}

@media only screen and (min-width: 1024px) {
    body {
        font-size: 16px;
        width: 75%;
    }
}`
