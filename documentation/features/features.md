# Features

All of allmarks' features in detail

This is an unordered list of the most prominent features of allmark:

1. Renders [GitHub Flavored MarkDown](https://help.github.com/articles/github-flavored-markdown/)
2. Full text search (+ Autocomplete)
3. Live-Reload / Live-Editing (via WebSockets)
4. Document Tagging
5. Tag Cloud
6. Documents By Tag
7. HTML Sitemap
8. XML Sitemap
9. robots.txt
10. RSS Feed
11. Print Preview
12. JSON Representation of Documents
13. Hierarchical Document Trees
14. Repository Navigation
	- Top-Level Navigation
	- Bread-Crumb Navigation
	- Previous and Next Items
	- Child-Documents
15. Image Thumbnails
16. Markdown Extensions
	- Image Galleries
	- File Preview
	- Displaying Folder Contents
	- Video Player Integration
	- Audio Player Integration
	- Repository cross-links by alias
17. Different Item Types (Repository, Document, Presentation)
18. Document Meta Data
	- Author
	- Tags
	- Document Alias
	- Creation Date
	- Last Modified Date
	- Language
	- Geo Location
19. Default Theme
	- Responsive Design
	- Lazy Loading for images and videos
	- Syntax Highlighting
20. Presentation Mode
21. Rich Text Conversion (Download documents as .rtf files)
22. Image Thumbnail Generation
23. HTTPS Support
	- Reference custom SSL certificates via `.allmark/config` from the `.allmark/certs` folder
	- Generates self-signed SSL certificates on-the-fly if no certificate is configured
24. Basic-Authentication
	- For an additional level of security allmark will only allow basic-authentication over SSL.
	- You can add users to the `.allmark/users.htpasswd` file using the tool [htpasswd](http://httpd.apache.org/docs/2.2/programs/htpasswd.html)
25. Parallel hosting of HTTP/HTTPS over IPv4 and/or IPv6
26. Short links: If you assign an alias to a document you can reach that document via short/direct link (e.g. `http://repo.com/!an-alias`). An overview of all available short links can be reached under `http://repo.com/!`.
27. You can use [Emojis](http://www.emoji-cheat-sheet.com/) in your markdown code :dancers:

---

created at: 2015-08-03
modified at: 2015-08-03
author: Andreas Koch
tags: Features, Documentation
alias: features
