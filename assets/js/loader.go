package jspacker

// This is our very simple module loader to bundle up JS files.
const moduleLoader = `"use strict";

const _modules = {
${PACK_MODULES}
};

function require(path) {
	const loader = _modules[path];
	
	if(!loader)
		throw "Module not found: " + path;
	
	if(loader.exports !== undefined)
		return loader.exports;
	
	loader.exports = {};
	loader(loader.exports);
	
	return loader.exports;
}`
