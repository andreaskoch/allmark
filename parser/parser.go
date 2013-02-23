package parser

func GetParser(lines []string) func() Document {

	return func() Document {
		return CreateDocument(lines)
	}

}
