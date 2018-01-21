var Username;
var lastBuy;
var lastSell;

$("#set-user").on("submit", function(e) {
	e.preventDefault();

	var formData = {
		'User'				 : $('input[name="User"').val()
	};
	console.log("Set");
	Username = $('input[name="User"').val();

	$.ajax({
		url:'/AddUser',
		type:'post',
		dataType: "json",
		contentType: 'application/json; charset=utf-8',
		data: JSON.stringify(formData),
		success:function(data){
			console.log(data);
			$('#current-user').replaceWith(document.createTextNode( "Logged in as : " + data.User ) );
		},
		error: function (xhr, status) {
			alert("Sorry, there was a problem!");
		},
	});
	return false;
});


$("#add-funds").on("submit", function(e) {
	e.preventDefault();

	var formData = {
		'UserId'		     : Username,
		'Amount'			 : $('input[name="AmountFunds"').val()
	};

	$.ajax({
		url:'/AddFunds',
		type:'post',
		dataType: "json",
		contentType: 'application/json; charset=utf-8',
		data: JSON.stringify(formData),
		success:function(data){
			console.log(data);
			$('#funds-added').replaceWith(document.createTextNode( "Added funds : " + data.Amount ) );
		},
		error: function (xhr, status) {
			alert("Sorry, there was a problem!");
		},
	});
	return false;
});


$("#get-quote").on("submit", function(e) {
	e.preventDefault();

	var formData = {
		'StockSymbol'            : $('input[name="Stock"').val(),
		'UserId'				 : Username
	};

	$.ajax({
		url:'/GetQuote',
		type:'post',
		dataType: "json",
		contentType: 'application/json; charset=utf-8',
		data: JSON.stringify(formData),
		success:function(data){
			console.log(data);
			$('#stock-price').html( data.Stock +":"+ data.Price );
		},
		error: function (xhr, status) {
			alert("Sorry, there was a problem!");
		},
	});
	return false;
});


$("#buy-stock").on("submit", function(e) {
	e.preventDefault();

	var formData = {
		'UserId'				 : Username,
		'StockSymbol'				 : $('input[name="StockBuy"').val(),
		'Amount'			 : Number($('input[name="AmountBuy"').val()),
	};

	$.ajax({
		url:'/BuyStock',
		type:'post',
		dataType: "json",
		contentType: 'application/json; charset=utf-8',
		data: JSON.stringify(formData),
		success:function(data){
			lastBuy = getTime();
			console.log(data);
			$('#confirm-buy').html( "Buy : " + data.Amount + "  Of : " + data.Stock );
			var r = confirm("Buy stock?");
				if (r == true) {
					if(!(getTime()>lastBuy+6000)){
						cancelorConfirm(data, '/CommitBuy', 'Commited');
					}
				} else {
					if(!(getTime()>lastBuy+6000)){
						cancelorConfirm(data, '/CancelBuy', 'Cancelled');
					}
				}
		},
		error: function (xhr, status) {
			alert("Sorry, there was a problem!");
		},
	});
	return false;
});

function cancelorConfirm(d, url, msg){
	var formData = {
		'UserId'				 : d.User
	};
	$.ajax({
		url: url,
		type:'post',
		dataType: "json",
		contentType: 'application/json; charset=utf-8',
		data: JSON.stringify(d),
		success:function(data){
			console.log(data);
			alert(msg);
		},
		error: function (xhr, status) {
			alert("Sorry, there was a problem!");
		},
	});
	return false;
}

$("#sell-stock").on("submit", function(e) {
	e.preventDefault();

	var formData = {
		'UserId'				 : Username,
		'StockSymbol'				 : $('input[name="StockSell"').val(),
		'Amount'			 : Number($('input[name="AmountSell"').val()),
	};

	$.ajax({
		url:'/SellStock',
		type:'post',
		dataType: "json",
		contentType: 'application/json; charset=utf-8',
		data: JSON.stringify(formData),
		success:function(data){
			console.log(data);
			lastSell = getTime();
			$('#confirm-sell').html( "Sell : " + data.Amount + "  Of : " + data.Stock );
			var r = confirm("Sell stock?");
				if (r == true) {
					if(!(getTime()>lastSell+6000)){
						cancelorConfirm(data, '/CommitSell', 'Commited');
					}
				} else {
					if(!(getTime()>lastSell+6000)){
						cancelorConfirm(data, '/CancelSell', 'Cancelled');
					}
				}
		},
		error: function (xhr, status) {
			alert("Sorry, there was a problem!");
		},
	});
	return false;
});


(function($) {

	skel.breakpoints({
		xlarge:	'(max-width: 1680px)',
		large:	'(max-width: 1280px)',
		medium:	'(max-width: 980px)',
		small:	'(max-width: 736px)',
		xsmall:	'(max-width: 480px)'
	});

	$(function() {

		var	$window = $(window),
			$body = $('body');

		// Disable animations/transitions until the page has loaded.
			$body.addClass('is-loading');

			$window.on('load', function() {
				window.setTimeout(function() {
					$body.removeClass('is-loading');
				}, 100);
			});

		// Fix: Placeholder polyfill.
			$('form').placeholder();

		// Prioritize "important" elements on medium.
			skel.on('+medium -medium', function() {
				$.prioritize(
					'.important\\28 medium\\29',
					skel.breakpoint('medium').active
				);
			});

	});

})(jQuery);