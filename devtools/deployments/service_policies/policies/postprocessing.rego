package postprocessing

import future.keywords.if
import data.utils

default granted := true

granted = false if {
    not utils.is_extension_allowed(input.resource.name)
}

granted = false if {
    bytes := quad4.resource.download(input.resource.url)
    mimetype := quad4.mimetype.detect(bytes)

    not utils.is_mimetype_allowed(mimetype)
}
