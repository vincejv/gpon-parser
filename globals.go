package main

var cachedPage = new(DocPage)         // main page
var cachedPage2 = new(DocPage)        // additional page
var cachedGponData = new(GponPayload) // zlt g3000a payload
var gponSvc OntDevice
