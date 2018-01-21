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
			var hasSubmenu = $(el).hasClass("has-submenu");
			$(global.menuClass + " .is-active").removeClass("is-active");
			$(el).addClass("is-active");

			// if (hasSubmenu) {
			// 	$(el).find("ul").slideDown();
			// }
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
});

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