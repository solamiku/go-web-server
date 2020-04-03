function preProcess(data) {
	var jdata = {};
	try {
		jdata = JSON.parse(data);
	} catch (e) {
		console.error("parse ", data, " failed");
		return null;
	}
	if (jdata.err) {
		thint.error(jdata.err);
		return null;
	}
	if (jdata.msg) {
		thint.success(jdata.msg);
	}
	return jdata;
};

var checkBoxModal;

$(function() {
	checkBoxModal = new CheckBoxList(".checkbox-modal", {});
});

function userclick(obj) {
	var self = $(obj);
	var attrNames = [];
	self.parent().parent().find(".authNames").each(function() {
		attrNames.push($(this).html());
	});

	var opts = [];
	self.find(".auth").each(function(idx) {
		opts.push({
			desc: attrNames[idx + 2],
			checked: $(this).html() == "true"
		});
	});

	checkBoxModal.active({
		title: "权限设置-%s".format(self.find("td:nth-child(2)").html()),
		opt: opts,
	}, function(datas) {
		console.log(datas);
		$.post("/setAuth", {
			auth: datas,
			uid: parseInt(self.find("td").first().html())
		}, function(ret) {
			var jdata = preProcess(ret);
			if (jdata) {
				self.find(".auth").each(function(idx) {
					$(this).html(jdata.auths[idx]);
				});
			}
		});
	});
};

/**
 * check box list
 * require bootstrap and jquery
 *
 * el - selecter
 * opt config
 * 	@clear - clear all el content or not.
 */
var CheckBoxList = function(el, opt) {
	this.obj = $(el);
	this._callback = function() {
		console.log("need to override")
	};

	if (this.obj.length == 0) {
		this.obj = $("<div></div>").appendTo($("body"));
	}

	var self = this;
	init(opt.clear);


	function init(clear) {
		if (clear) {
			self.obj.html("");
		}
		self.obj.addClass("modal fade");
		self.obj.attr({
			"tabindex": "-1",
			"role": "dialog",
			"aria-hidden": "true"
		});
		var dialog = $(`<div class="modal-dialog"></div>`).appendTo(self.obj);
		var content = $(`<div class="modal-content"></div>`).appendTo(dialog);
		self.header = $(`<div class="modal-header">
							<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
                    		<h4 class="modal-title"></h4>
                    	</div>`).appendTo(content);
		self.body = $(`<div class="modal-body" style="margin:10px;max-height:500px;overflow: auto;">
                	 </div>`).appendTo(content);
		var footer = $(`<div class="modal-footer" id="input_modal_footer"></div>`).appendTo(content);
		var btn = $(`<button class="btn btn-default">确定</button>`).appendTo(footer);
		btn.click(function() {
			var datas = [];
			self.body.find(".custom-control-input").each(function() {
				datas.push($(this).is(":checked"));
			});
			self._callback(datas);
			self.obj.modal('toggle')
		});
	};
};

CheckBoxList.prototype.active = function(datas, callback) {
	this._callback = callback;
	this.header.find(".modal-title").html(datas.title);
	this.body.html("");
	for (var i in datas.opt) {
		var checkbox = $(`<div class="row" style="margin:10px">
							<div class="col-sm-2"></div>
							<div class="col-sm-4"><b>%s:</b></div>
							<div class="col-sm-4">
								<label class="custom-control custom-checkbox">
					        	<input type="checkbox" class="custom-control-input">
					        	<span class="custom-control-indicator"></span>
				           		</label>
				            </div>
				        </div>`.format(datas.opt[i].desc)).appendTo(this.body);
		if (datas.opt[i].checked) {
			checkbox.find("input").attr("checked", "checked");
		}
	}
	this.obj.modal('toggle');
};

/*
 
 */