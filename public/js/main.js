var map = L.map('map').setView([39.2847064,-76.620486], 11);
var markers = new L.MarkerClusterGroup({disableClusteringAtZoom: 10,});

function mapItems() {
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

function highlight(id) {
	$('.item').removeClass('highlight');
	$('#item'+id).addClass('highlight');
}

$(document).ready(function() {
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

	mapItems();
});


