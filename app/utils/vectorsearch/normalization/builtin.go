package normalization

import (
	_ "nova-factory-server/app/utils/vectorsearch/normalization/lowercase"
	_ "nova-factory-server/app/utils/vectorsearch/normalization/regex"
	_ "nova-factory-server/app/utils/vectorsearch/normalization/replace"
	_ "nova-factory-server/app/utils/vectorsearch/normalization/shopcategory"
	_ "nova-factory-server/app/utils/vectorsearch/normalization/whitespace"
)
