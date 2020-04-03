/*for animate.css animation end callback*/
$.fn.extend({
	animateCss: function(animationName, callback) {
		var animationEnd = 'webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend';
		this.addClass('animated ' + animationName).one(animationEnd, function() {
			$(this).removeClass('animated ' + animationName);
			if (callback) {
				callback($(this));
			}
		});
		return this;
	}
});

$(function() {
	var Dashboard = function() {
		var global = {
			tooltipOptions: {
				placement: "right"
			},
			menuClass: ".c-menu"
		};

		var menuChangeActive = function menuChangeActive(el) {
			el = $(el);
			var hasSubmenu = el.hasClass("has-submenu");
			$(global.menuClass + " .is-active").removeClass("is-active");
			el.addClass("is-active");

			// if (hasSubmenu) {
			// 	$(el).find("ul").slideDown();
			// }

			//show tab
			var cur = $(".content-tab-active");
			var to = $(el.attr("href"));
			if (cur == to) {
				return;
			}
			cur.removeClass("content-tab-active");
			cur.addClass("content-tab");
			to.removeClass("content-tab");
			to.addClass("content-tab-active");
			to.addClass("animated fadeIn");
		};

		var sidebarChangeWidth = function sidebarChangeWidth() {
			var $menuItemsTitle = $("li .menu-item__title");

			$("body").toggleClass("sidebar-is-reduced sidebar-is-expanded");
			$(".hamburger-toggle").toggleClass("is-opened");

			if ($("body").hasClass("sidebar-is-expanded")) {
				$('[data-toggle="tooltip"]').tooltip("destroy");
			} else {
				$('[data-toggle="tooltip"]').tooltip(global.tooltipOptions);
			}
		};

		return {
			init: function init() {
				$(".js-hamburger").on("click", sidebarChangeWidth);

				$(".js-menu li").on("click", function(e) {
					menuChangeActive(e.currentTarget);
				});

				$('[data-toggle="tooltip"]').tooltip(global.tooltipOptions);

				RingNotice = new RingNotices($(".ringnotice"));
			}
		};
	}();

	Dashboard.init();

	$("ul.loginbox").on("click", function(e) {
		e.stopPropagation();
	});

	/*panel*/
	$(".flip").on("click", function(e) {
		$(".flip").each(function() {
			var panel = $(this).parent().parent().find(".objpanel");
			if (this == e.currentTarget) {
				panel.slideToggle("quick");
			} else {
				panel.hide("quick");
			}
		});
	});

	$(".btn-command").click(function() {
		var self = $(this);
		var args = [];
		var server = self.attr("server");
		args.push("server=" + server);
		var command = self.attr("command");
		args.push("command=" + command);

		//parse clientp
		var clientpStr = self.attr("clientp");
		var params = [];
		if (clientpStr) {
			var clientps = clientpStr.split(";");
			for (var i in clientps) {
				var clientp = clientps[i];
				var tag2id = clientp.split(":");
				if (tag2id.length == 2) {
					params.push({
						type: "input",
						t: tag2id[0],
						v: "",
						id: tag2id[1],
					});
				}
			}
		}

		if (params.length > 0) {

			inputModalbox({
				title: "相关参数",
				inputs: params,
			}, function(datas) {
				for (var i in datas) {
					var val = datas[i];
					var pname = params[i].id;
					args.push(pname + "=" + val);
				}
				startIframe();
			});
		} else {
			startIframe();
		}

		//send server command
		function startIframe() {
			var tConsole = self.parent().parent().find(".command-ret");
			tConsole.html("");
			var href = "/command?" + args.join("&");
			var iframe = $(`<iframe  id="MainFrame"  frameborder="no" scrolling="no" border="0" src="%s" width="100%">`.format(href));

			iframe.appendTo(tConsole);
		}
	});

	//ifream自动适应高度,div滚动到最下方
	setInterval(function() {
		var iframe = $("iframe");
		if (!iframe || iframe.length == 0) {
			return;
		}
		var idx = 0;
		iframe.each(function(idx) {
			var iframeHeight = 0;
			if (navigator.userAgent.indexOf("Firefox") > 0) { // Mozilla, Safari, ...  
				iframeHeight = this.contentDocument.body.offsetHeight;
			} else if (navigator.userAgent.indexOf("MSIE") > 0) { // IE  
				iframeHeight = MainFrame.document.body.scrollHeight; //IE这里要用MainFrame，不能用obj，切记  
			} else {
				iframeHeight = $(this).contents().height();
			}
			if (LastIframeHeight[idx] != iframeHeight) {
				LastIframeHeight[idx] = iframeHeight;
				$(this).css("height", iframeHeight);
				var scrollHeight = $(this).parent().prop("scrollHeight");
				$(this).parent().scrollTop(scrollHeight, 20);
			}
		});

	}, 100);

	let start = $(".logstart");
	let end = $(".logend");
	if (start.length > 0 && end.length > 0) {
		let now = new Date();
		let hourpre = new Date(now.getTime() - 3600 * 1000);
		start.val(hourpre.format("%Y-%m-%dT%H:%M"));
		end.val(now.format("%Y-%m-%dT%H:%M"))
	}
	refreshAllDownList();
	refreshAllErrList();
});

function leitinguidlog(obj) {
	let self = $(obj);
	let uidinput = self.parent().find(".loguid");
	let uid = parseInt(uidinput.val());
	if (isNaN(uid)) {
		thint.error("uid异常");
		return;
	}
	let aidinput = self.parent().find(".logaid");
	let aid = parseInt(aidinput.val());
	if (isNaN(aid)) {
		thint.error("aid异常");
		return;
	}
	let start = self.parent().find(".logstart").val();
	let end = self.parent().find(".logend").val();
	let processTimestr = function(t) {
		t = t.replace("T", " ");
		t += ":00+08:00";
		return t
	}
	//	console.log(uid, start, end, processTimestr(start), processTimestr(end));
	let logdownlist = self.parent().parent();

	$.post("/leitinglogcmd", {
		cmd: "extract",
		logid: logdownlist.parent().attr("logid"),
		uid: parseInt(uid),
		aid: parseInt(aid),
		start: processTimestr(start),
		end: processTimestr(end),
	}, function(data) {
		let jdata = ToJson(data);
		console.log(jdata);
		refreshDownList(logdownlist);
	});
}

function ToJson(jsonstr) {
	try {
		let d = JSON.parse(jsonstr);
		if (d.err) {
			thint.error(d.err);
		}
		return d;
	} catch (e) {
		console.log("try parse ", jsonstr, e);
		return {}
	}
}

function refreshAllDownList() {
	let downlist = $(".logdownlist");
	if (downlist.length == 0) {
		return;
	}
	downlist.each(function() {
		let self = $(this);
		refreshDownList($(this));
	})
}

function quickRefreshDownList(obj) {
	let self = $(obj);
	let downlist = self.parent().parent();
	refreshDownList(downlist);
}

function quickRefreshErrList(obj) {
	let self = $(obj);
	let errlist = self.parent().parent();
	refreshErrList(1, errlist);
}

function refreshDownList(obj) {
	let self = $(obj);
	let logid = self.parent().attr("logid");
	let tbody = self.find("tbody");
	$.post("/leitinglogcmd", {
		logid: logid,
		cmd: "downlist"
	}, function(data) {
		tbody.html("");
		console.log(data);
		var jdata = {};
		try {
			jdata = JSON.parse(data)
		} catch (e) {}
		let list = jdata.list;
		list.sort(function(a, b) {
			return a.Start - b.Start;
		});
		for (let i in list) {
			let opt = list[i];
			let tr = $("<tr></tr>").appendTo(tbody);
			$(`<td>%s</td>`.format(opt.File)).appendTo(tr);
			$(`<td>%s</td>`.format(calsize(opt.Size))).appendTo(tr);
			$(`<td>%s</td>`.format(opt.Ip)).appendTo(tr);
			$(`<td>%s</td>`.format(new Date(opt.Start * 1000).format())).appendTo(tr);
			let end = "";
			if (opt.Done == 1 || opt.Done == -1) {
				end = new Date(opt.End * 1000).format();
			}
			$(`<td>%s</td>`.format(end)).appendTo(tr);
			let st = ""
			switch (opt.Done) {
				case 0:
					st = "正在下载";
					break;
				case 1:
					let path = opt.Path.replace("serverfs/", "") + "/" + opt.File
					st = `<a onclick="downfile('%s', '%s')">下载到本地</a>`.format(opt.File, path);
					st += ` | <a onclick="openfile('%s')">直接打开</a>`.format(path);
					break;
				default:
					st = "异常";
					break;
			}
			$(`<td>%s</td>`.format(st)).appendTo(tr);
			$(`<td>%s</td>`.format(opt.Err)).appendTo(tr);
		}
	});
}

function downfile(file, path) {
	let a = $('<a href="%s" download="%s">测试</a>'.format(path, file)).appendTo($("body"));
	a.get(0).click();
	a.remove();
}

function openfile(path) {
	window.open(path);
}

function calsize(size) {
	let nSize = parseInt(size);
	if (isNaN(nSize)) {
		return "0B"
	}
	if (nSize > 1024 * 1024) {
		return (nSize / 1024 / 1024).toFixed(2) + "MB"
	}
	if (nSize > 1024) {
		return (nSize / 1024).toFixed(2) + "KB"
	}
	return "0B"
}

function refreshAllErrList() {
	let errlist = $(".logerrlist");
	if (errlist.length == 0) {
		return;
	}
	errlist.each(function() {
		let self = $(this);
		refreshErrList(0, self);
	});
}

function refreshErrList(force, obj) {
	let self = $(obj);
	let logid = self.parent().attr("logid");
	let tbody = self.find("tbody");
	console.log(self, logid, force);
	$.post("/leitinglogcmd", {
		cmd: "errlist",
		force: force,
		logid: logid,
	}, function(data) {
		tbody.html("");
		console.log(data);
		let a = self.find("a");
		var jdata = {};
		try {
			jdata = JSON.parse(data)
		} catch (e) {}
		a.html("上次刷新时间:%d, 当前:%s.点击强制刷新".format(jdata.last, jdata.msg));
		let list = jdata.list;
		list.sort(function(a, b) {
			return a.File - b.File;
		});
		for (let i in list) {
			let opt = list[i];
			let tr = $("<tr></tr>").appendTo(tbody);
			$(`<td>%s</td>`.format(opt.File)).appendTo(tr);
			$(`<td>%s</td>`.format(calsize(opt.Size))).appendTo(tr);
			$(`<td>%s</td>`.format(opt.Mod)).appendTo(tr);
			$(`<td><a onclick="openfile('%s')">直接打开</a></td>`.format(opt.Path.replace("serverfs/", "") + opt.File)).appendTo(tr);
		}
	});
}


var LastIframeHeight = {};

var RingNotice;

function RingNotices(el) {
	this._notices = [];
	this._elNotice = null;
	this._elNoticeClass = "c-badge c-badge--header-icon animated shake";
	this._el = el;

	function init() {
		el.tooltip();
	};
	init();
};

RingNotices.prototype.addNotice = function(str) {
	this._notices.push(str);
	var l = this._notices.length;
	if (this._elNotice) {
		this._elNotice.remove();
	}
	this._elNotice = $("<span></span>").appendTo(this._el);
	this._elNotice.addClass(this._elNoticeClass);
	this._elNotice.html(l);
};

function jumpTo(url) {
	window.open(url);
};

// dbflush related
// 
function dbflush(obj) {
	let self = $(obj);
	let tmpl = self.parent().find('select[name="tmpl"]');
	let tmplselect = tmpl.find('option:selected').text();
	let db = self.parent().find('select[name="db"]');
	let dbselect = db.find('option:selected').text();
	let uid = parseInt(self.parent().find('input[name="uid"]').val());
	let aid = parseInt(self.parent().find('input[name="aid"]').val());
	// console.log(tmpl, db, uid);
	if (isNaN(uid) || isNaN(aid) || uid <= 0 || aid <= 0) {
		thint.error('uid或者aid异常');
		return;
	}
	confirmbox({
		title: '是否覆盖存档',
		content: `是否使用<b>${tmplselect}</b>覆盖到数据库:<b>${dbselect}</b>的目标<b>aid:${aid},uid:${uid}</b>上?`
	}, () => {
		self.html('<i class="fa fa-spinner fa-spin"></i>执行中');
		self.attr('disabled', true);
		$.post('/dbflush', {
			cmd: 'exec',
			tmpl: tmpl.val(),
			db: db.val(),
			uid: uid,
			aid: aid,
		}, (data) => {
			let jdata = JSON.parse(data);
			if (jdata.err) {
				thint.error(jdata.err, {
					"timeOut": "5000"
				});
			}
			if (jdata.affected) {
				thint.info(`执行完毕，影响行数:${jdata.affected},错误:${jdata.err}`, {
					"timeOut": "5000"
				});
			}
			self.html('覆盖到数据库');
			self.attr('disabled', false);
		});
	});
}

function newtmpl() {
	tmplmodify(0, '', '', '');
}

function modifytmpl(obj) {
	let self = $(obj);
	let childs = self.children();
	let id = childs.eq(0);
	let info = childs.eq(1);
	let tmpl = childs.eq(2);
	let auth = childs.eq(3);
	tmplmodify(parseInt(id.text()), info.text(), tmpl.text(), auth.text());
}

function tmplmodify(id, info, tmpl, auths) {
	let inputs = [];
	inputs.push({
		type: 'input',
		t: 'id',
		v: id,
		d: true,
	});
	inputs.push({
		type: 'input',
		t: '模板说明',
		v: info,
	});
	inputs.push({
		type: 'textarea',
		t: '具体内容',
		v: tmpl,
	});
	inputs.push({
		type: 'input',
		t: '权限',
		v: auths,
	});
	inputModalbox({
		inputs: inputs
	}, (datas) => {
		$.post("/dbflush", {
			cmd: 'tmpl',
			id: parseInt(datas[0]),
			info: datas[1],
			tmpl: datas[2],
			auth: datas[3]
		}, (data) => {
			let jdata = JSON.parse(data);
			if (jdata.err) {
				thint.error(jdata.err, {
					"timeOut": "5000"
				});
			}
			if (jdata.info) {
				thint.info(jdata.info);
			}
		})
	});
}

function newcfg() {
	cfgmodify(0, '', '', '');
}

function modifycfg(obj) {
	let self = $(obj);
	let childs = self.children();
	let id = childs.eq(0);
	let info = childs.eq(1);
	let dest = childs.eq(2);
	cfgmodify(parseInt(id.text()), info.text(), dest.text());
}

function cfgmodify(id, info, dest) {
	let inputs = [];
	inputs.push({
		type: 'input',
		t: 'id',
		v: id,
		d: true,
	});
	inputs.push({
		type: 'input',
		t: '配置说明',
		v: info,
	});
	inputs.push({
		type: 'input',
		t: '目标数据库',
		v: dest,
	});
	inputModalbox({
		inputs: inputs
	}, (datas) => {
		$.post("/dbflush", {
			cmd: 'cfg',
			id: parseInt(datas[0]),
			info: datas[1],
			dest: datas[2],
		}, (data) => {
			let jdata = JSON.parse(data);
			if (jdata.err) {
				thint.error(jdata.err, {
					"timeOut": "5000"
				});
			}
			if (jdata.info) {
				thint.info(jdata.info);
			}
		})
	});
}