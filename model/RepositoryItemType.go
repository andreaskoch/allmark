package model

type RepositoryItemType byte

const (
	NONE = iota
	REPOSITORYDESCRIPTION
	DOCUMENT
	COMMENT
	LOCATION
	MESSAGE
	TAG
)

func (itemType RepositoryItemType) String() string {
	switch itemType {

	case REPOSITORYDESCRIPTION:
		{
			return "Repository Description"
		}

	case DOCUMENT:
		{
			return "Document"
		}

	case COMMENT:
		{
			return "Comment"
		}

	case LOCATION:
		{
			return "Location"
		}

	case MESSAGE:
		{
			return "Message"
		}

	case TAG:
		{
			return "Tag"
		}

	}

	return "Unidentified Repository Item Type"
}
