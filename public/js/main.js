function dispatch(type, data) {
	var jdata = {};
	try {
		jdata = JSON.parse(data);
	} catch (e) {
		console.error("parse ", data, " failed");
	}
	if (jdata.err) {
		thint.error(jdata.err);
		return;
	}
	if (jdata.msg) {
		thint.success(jdata.msg);
	}
	switch (type) {}
};