// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themes

	const screenCss = `
html {
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
    margin: 10px auto;
    padding: 1em;
    outline: 1300px solid #FAFAFA;
}

body>nav {
    font-size: 0.8em;
}

body>nav>ul.breadcrumb {
    list-style: none;
    margin: 0;
    padding: 0;
}

body>nav>ul.breadcrumb>li {
    display: inline;
}

body>nav>ul.breadcrumb>li:after {
    content: ">";
}

body>nav>ul.breadcrumb>li:last-child:after {
    content: "";
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
    margin: 0 0 15px 0;
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
    margin: 1em 0;
    padding: 0 0 0 2em;
}

li p:last-child {
    margin: 0;
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

article>.description {
    font-size: 1.2em;
}

.subentries {
    list-style: none;
    padding: 5px 0 5px 0;
    margin: 0 0 0 15px;
}

.subentries>.subentry {
    margin: 0 0 15px 0;
}

.subentries>.subentry:nth-child(odd) {
    background-color:#eee;
}

.subentries>.subentry:nth-child(even) {
    background-color:transparent;
}

.imagegallery>h1 {
    font-size: 1.2em;
}

.imagegallery ol {
    list-style: none;
    margin-left: 0;
}

.filelinks>h1 {
    font-size: 1.2em;
}

.filelinks ol {
    margin-left: 0;
}

.collection>h1 {
    font-size: 1.2em;
}

.csv {
    margin: 20px 0 0 20px;
    overflow: auto;
}

.csv>h1 {
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
}

.pdf>h1 {
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

.video>h1 {
    font-size: 1.2em;
}

.audio>h1 {
    font-size: 1.2em;
}

.presentation nav {
    float: left;
    margin: 25px 0 15px 0;
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