var map = L.map('map').setView([39.2847064,-76.620486], 11);

function mapItems() {
	$('.item').each(function(i) {
		elm = $(this);
		// add a marker in the given location, attach some popup content to it and open the popup
		L.marker([elm.data('lat'),elm.data('long')]).addTo(map)
		    .bindPopup(elm.children('.title').text())
		    .openPopup();
		console.log(elm.children('.title'));
	});
}

$(window).load(function() {
	// add an OpenStreetMap tile layer
	//L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {

	// add a MapQuest tile layer
	L.tileLayer('http://otile3.mqcdn.com/tiles/1.0.0/osm/{z}/{x}/{y}.jpg', {
	    attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
	}).addTo(map);

	mapItems();
});


