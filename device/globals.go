package device

import "github.com/vincejv/gpon-parser/util"

var cachedPage = new(util.DocPage)    // main page
var cachedPage2 = new(util.DocPage)   // additional page
var cachedGponData = new(GponPayload) // zlt g3000a payload
var GponSvc OntDevice
