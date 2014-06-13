function mapItems(markers) {
	$('.item').each(function(i) {
		elm = $(this);
		// add a marker in the given location, attach some popup content to it and open the popup
		marker = L.marker([elm.data('lat'),elm.data('long')],{
			item_id: elm.data('id'),
			title: elm.data('id'),
			icon: L.AwesomeMarkers.icon({
				prefix: 'fa',
				icon: elm.data('icon'),
				markerColor: elm.data('color'),
			}),
		})

		markers.addLayer(marker);

	});
}

function reply(id) {
	$('#tools'+id).hide();
	$('#reply'+id).show();
}

function vote(t, id, v) {
	elm = $('#'+t+id)[0];
	$(elm).find('.vote').removeClass('voted');
	if(v>0) {
		$(elm).find('.upvote').addClass('voted');
	}
	if(v<0) {
		$(elm).find('.downvote').addClass('voted');
	}
}


function highlight(id) {
	$('.item').removeClass('highlight');
	$('#i'+id).addClass('highlight');
}

$(document).ready(function() {
	$(".commentform").submit(function (event) {
		event.preventDefault();

		var url = $(this).attr('action');
		var comments = $(this).siblings(".comments");

		$.ajax({
		   type: "POST",
		   url: url,
		   data: $(this).serialize(),
		   success: function(data) {
			comments.prepend(data);
		   }
		 });

	});

	$(".reply").click(function (event) {
		event.preventDefault();

		id = $(this).data('id');
		reply(id);
	});	

	$(".score .vote").click(function (event) {
		event.preventDefault();
		elm = $(this)
		overall = elm.siblings('.overall')[0];
		var preSUV = parseInt($(overall).data('suv'));
		var baseS = parseInt(overall.innerHTML) - preSUV; 

		$.ajax({
			url: '/vote',
			method: 'POST',
			data: {
				'pt': elm.data('pt'),
				'p': elm.data('p'),
				'v': elm.data('v')
			},
			dataType: "json",
			success: function(data) {
				overall.innerHTML = baseS + data.v;
				$(overall).data('suv', data.v);
				vote(data.pt, data.p, data.v);
			}
		});
	});
	
	if($('#map').length) {
		var map = L.map('map').setView([39.2847064,-76.620486], 11);
		var markers = new L.MarkerClusterGroup({disableClusteringAtZoom: 10,});

		// add an OpenStreetMap tile layer
		//L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {

		// add a MapQuest tile layer
		L.tileLayer('http://otile{s}.mqcdn.com/tiles/1.0.0/osm/{z}/{x}/{y}.jpg', {
		    subdomains: '1234',
		    attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors',
		}).addTo(map);

		map.addLayer(markers);

		/*map.on('popupopen', function(e) {
			console.log(e); // e is an event object (MouseEvent in this case)
		});*/

		markers.on('click', function (d) {
			var item_id = d.layer.options.item_id;
			
			if ( $('#item'+item_id) ){
				highlight(item_id);
			}
		});

		mapItems(markers);
	}
});


