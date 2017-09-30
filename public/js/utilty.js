/**
 *   通用工具库
 *
 * > string格式化
 * > thint 基于 toastr库的一些简单调用封装
 * > confirmbox 模态式确认框
 * > inputModalbox 模态式输入框 需要基于bootstrap-select库
 * > 时间格式化
 * > 带超时的ajax请求
 *
 *
 *
 *
 *
 * 
 *   by lipm
 */

/**
 * String对象的格式化功能
 */
String.prototype.format = function() {
	var args = arguments
	str = this;
	for (var x in args) {
		str = str.replace(/\%(\w)/, args[x])
	}
	return str
};



/**
 * toastr 提示相关
 */
var thint, Retpro;
if (toastr) {
	toastr.options = {
		"closeButton": false,
		"debug": false,
		"newestOnTop": false,
		"progressBar": false,
		"positionClass": "toast-top-center",
		"preventDuplicates": false,
		"onclick": null,
		"showDuration": "500",
		"hideDuration": "1000",
		"timeOut": "2000",
		"extendedTimeOut": "1000",
		"showEasing": "swing",
		"hideEasing": "linear",
		"showMethod": "fadeIn",
		"hideMethod": "fadeOut"
	}


	//alert hint
	var thint = {
		max: 1,
		cnt: [0, 0, 0, 0],
		error: function(text, cfg) {
			this.check(1);
			toastr.error(text, "", cfg);
		},
		info: function(text, cfg) {
			this.check(2);
			toastr.info(text, "", cfg);
		},
		warn: function(text, cfg) {
			this.check(3);
			toastr.warning(text, "", cfg);
		},
		success: function(text, cfg) {
			this.check(4);
			toastr.success(text, "", cfg);
		},
		check: function(idx) {
			if (this.cnt[idx] > this.max) {
				this.cnt[idx] = 1;
				toastr.remove();
			} else {
				this.cnt[idx]++;
			}
		}
	};

	var Retpro = {
		normal: function(data) {
			if (data == "ok") {
				location.reload(true);
				return;
			}
			thint.info(data, {
				timeOut: 1000
			});
		}
	};
}


/**
 * confirm box
 */

function confirmbox(boxinfo, f) {
	var confirm = $("#confirm_modal");
	if (confirm.length <= 0) {
		$("body").append(`
		<div class="modal in" style="display: block;" id="confirm_modal">
			<div class="modal-dialog">
				<div class="modal-content">
					<div class="modal-header">
						<h4 class="modal-title" id="confirm_modal_title">Are you sure???</h4>
					</div>
					<div class="modal-body">
						<p id="confirm_modal_content">Are you sure you want to delete (this)?</p>
						<div class="row">
						  <div class="col-12-xs text-center">
								<button class="btn btn-success btn-md" id="confirm_modal_yes">Yes</button>
								<button class="btn btn-danger btn-md" id="confirm_modal_no">No</button>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
		`);
	}
	var confirm = $("#confirm_modal");
	if (boxinfo.title) {
		$("#confirm_modal_title").html(boxinfo.title);
	}
	if (boxinfo.content) {
		$("#confirm_modal_content").html(boxinfo.content);
	}
	if (boxinfo.btnyes) {
		$("#confirm_modal_yes").html(boxinfo.btnyes);
	}
	if (boxinfo.btnno) {
		$("#confirm_modal_no").html(boxinfo.btnno);
	}
	//每次都要重新绑定以便调用正确的回调函数
	//防止重复绑定点击时触发多次，需要先unbind
	$("#confirm_modal_yes").unbind();
	$("#confirm_modal_yes").click(function() {
		if (f) {
			f();
		}
		confirm.hide();
	});
	$("#confirm_modal_no").unbind();
	$("#confirm_modal_no").click(function() {
		confirm.hide();
	});
	confirm.show();
};


/**
 * input modal box
 */
function inputModalbox(args, callback, onSelect) {
	var args = args || {};
	var inputmodal = $("#input_modal");
	if (inputmodal.length <= 0) {
		inputmodal = $(`
			<div class="modal fade" id="input_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
				<div class="modal-dialog">
					<div class="modal-content">
						<div class="modal-header">
							<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
							<h3 class="modal-title" id="input_modal_title">标题</h3>
						</div>
						<div class="modal-body">
							<form role="form" id="input_modal_form">
							</form>
						</div>
						<div class="modal-footer" id="input_modal_footer">
							<button type="button" class="btn btn-default" opt="0">确定</button>
						</div>
					</div>
				</div>
			</div>
		`).appendTo($("body"));
	}
	var footer = inputmodal.find("#input_modal_footer");
	footer.html("");
	if (args.btn) {
		for (var i in args.btn) {
			var btn = args.btn[i];
			$(`<button type="button" class="btn btn-default" opt="%s">%s</button>`.format(i, btn)).appendTo(footer);
		}
	} else {
		$(`<button type="button" class="btn btn-default" opt="0">确定</button>`).appendTo(footer);
	}
	//每次都要重新绑定以便调用正确的回调函数
	//防止重复绑定点击时触发多次，需要先unbind
	footer.children("button").unbind();
	footer.children("button").click(function() {
		var childs = $("#input_modal_form").children();
		var datas = [];
		for (var i = 0; i < childs.length; i++) {
			var input = $(childs[i]);
			var obj = input.find("input[class!='form-control'],select,textarea");
			var val;
			switch (obj.attr("type")) {
				case "checkbox":
					val = obj.parent().hasClass("active");
					break;
				default:
					val = obj.val();
			}
			datas.push(val);
		}
		callback(datas, parseInt($(this).attr("opt")));
		inputmodal.modal("hide");
	});

	//init
	var form = inputmodal.find("#input_modal_form");
	form.html("");
	//title
	if (args.title) {
		inputmodal.find("#input_modal_title").html(args.title);
	}

	//inputs
	var focusObj;
	if (args.inputs) {
		for (var i in args.inputs) {
			var ops = args.inputs[i];
			switch (ops.type) {
				case "select":
					var op = $("<div class=\"row\"><div class=\"col-xs-1\"></div></div>").appendTo(form);
					var label = $("<div class=\"col-xs-3\"><label>%v : </label></div>".format(ops.t)).appendTo(op);
					var sdiv = $("<div class=\"col-xs-6\"><select class=\"selectpicker\" data-live-search=\"true\"></select></div>").appendTo(op);
					var select = sdiv.children("select");
					for (var i in ops.v) {
						var o = $("<option value=\"%v\" data-subtext=\"%v\">%v</option>".format(ops.v[i].id, ops.v[i].id, ops.v[i].val)).appendTo(select);
						if (ops.v[i].id == ops.s) {
							o.attr("selected", "selected");
						}
					}
					break;
				case "checkbox":
					var op = $("<div class=\"row\"><div class=\"col-xs-1\"></div></div>").appendTo(form);
					var label = $("<div class=\"col-xs-3\"><label>%v : </label></div>".format(ops.t)).appendTo(op);
					var input = $(`
						<div class=\"col-xs-6 btn-group\" data-toggle=\"buttons\">
							<label class=\"btn btn-default %v\">
								<input type=\"checkbox\" autocomplete=\"off\">
								<span class=\"glyphicon glyphicon-ok\"></span>
							</label>
						</div>`.format(ops.v == "true" ? "active" : "")).appendTo(op);
					break;
				case "textarea":
					var op = $("<div class=\"row\"></div>").appendTo(form);
					var textarea = $("<textarea style=\"margin-left:20px\"></textarea>").appendTo(op);
					var size = {
						w: 550,
						h: 400
					}
					if (ops.size) {
						size = ops.size;
					}
					textarea.height(size.h);
					textarea.width(size.w);
					if (ops.v) {
						textarea.val(ops.v);
					}
					if (!focusObj) {
						focusObj = textarea;
					}
					break;
				case "input":
				default:
					var op = $("<div class=\"row\" style=\"margin-top:10px\"><div class=\"col-xs-1\"></div></div>").appendTo(form);
					var label = $("<div class=\"col-xs-3\"><label>%v : </label></div>".format(ops.t)).appendTo(op);
					var input = $("<div class=\"col-xs-6\"><input type=\"text\" class=\"form-control input-sm\"/></div>").appendTo(op);
					if (ops.v) {
						input.children("input").val(ops.v);
					}
					if (ops.p) {
						input.children("input").attr("placeholder", ops.p);
					}
					if (!focusObj) {
						focusObj = input.children("input");
					}
					break;
			}
		}
	}

	inputmodal.unbind('shown.bs.modal');
	inputmodal.on('shown.bs.modal', function() {
		if (focusObj) {
			focusObj.focus();
			focusObj.select();
		}
	});
	inputmodal.modal("show");
	inputmodal.modal("show");

	$('.selectpicker').selectpicker({
		size: 20,
	});
	$('div.bootstrap-select').each(function() {
		var ul = $(this).children(".dropdown-menu").children("ul");
		ul.children("li").each(function() {
			var a = $(this).children("a");
			if (a.length > 0) {
				var text = a.children(".text").text();
				a.attr("data-tokens", H2P(text));
			}
		})
	});
	return inputmodal;
};


/**
 * ajax
 */

if (!$) {
	$ = {};
}
$.spost = function(url, data, f) {
	var args = {
		error: function(XMLHttpRequest, textStatus, errorThrown) {
			thint.error("请求发生错误");
		},
		success: function(data, textStatus) {
			// console.log(data, textStatus);
			f(data);
		},
		complete: function(XMLHttpRequest, textStatus) {
			// console.log(XMLHttpRequest, textStatus);
		},
		statusCode: {
			404: function() {
				thint.error("404, 页面不存在.");
			}
		},
		type: "POST",
		data: data
	}
	$.ajax(url, args);
};

/** 简单的get跳转 */
$.sget = function(url, args) {
	var str = url;
	str += "?";
	for (var i in args) {
		str += (i + "=" + args[i] + "&");
	}
	window.location.href = str;
};

/**
 *  time-related
 */
Date.prototype.format = function(format) {
	if (!format) {
		format = '%Y-%m-%d %H:%M:%S';
	}
	var self = this;
	return format.replace(/%[YmdHMS]/g, function(m) {
		switch (m) {
			case '%Y':
				return self.getFullYear();
			case '%m':
				m = 1 + self.getMonth();
				break;
			case '%d':
				m = self.getDate();
				break;
			case '%H':
				m = self.getHours();
				break;
			case '%M':
				m = self.getMinutes();
				break;
			case '%S':
				m = self.getSeconds();
				break;
			default:
				return m.slice(1);
		}
		return ('0' + m).slice(-2);
	});
};

Date.prototype.addDay = function(del) {
	var t = this.getTime();
	t += (del * 3600 * 1000 * 24);
	return new Date(t);
};


function getTimeZero(s) {
	var now;
	if (s) {
		now = new Date(s);
	} else {
		now = new Date();
	}
	return new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0, 0);
};