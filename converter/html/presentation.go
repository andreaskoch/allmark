// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"fmt"
	"strings"
)

func renderPresentation(html string) string {
	slides := strings.Split(html, "<hr />")
	presentationCode := fmt.Sprintf(`<section class="slide">%s</section>`, strings.Join(slides, `</section><section class="slide">`))
	return presentationCode
}
