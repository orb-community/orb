import { Component, Input, OnChanges } from '@angular/core';
import { MessagesService } from 'app/common/services/messages/messages.service';
import * as L from 'leaflet';
import { MqttService, IMqttMessage } from 'ngx-mqtt';
import { Gateway } from 'app/common/interfaces/gateway.interface';


@Component({
  selector: 'ngx-map-leaflet',
  templateUrl: './map.leaflet.component.html',
  styleUrls: ['./map.leaflet.component.scss'],
})
export class MapComponent implements OnChanges {
  map: L.Map;
  markersGroup = new L.LayerGroup();

  options = {
    layers: [
      L.tileLayer('https://api.mapbox.com/styles/v1/sasamainflux/cjz9wtqyr03711cp4bxew2f1k/tiles/256/{z}/{x}/{y}?' +
      'access_token=pk.eyJ1Ijoic2FzYW1haW5mbHV4IiwiYSI6ImNqejl3cGppODAybXAzbXFzcmcxZDE1cnEifQ.1E7DVFz5JPFiqnpP4GFvOA',
      {}),
    ],
    center: L.latLng({ lat: 45, lng: 20 }),
    zoom: 5,
    minZoom: 2,
    maxBounds: L.latLngBounds(L.latLng(-90, -180),  // southWest
      L.latLng(90, 180)),   // northEast
    maxBoundsViscosity: 1,
  };
  drawOptions = {
    position: 'topright',
    draw: {
      marker: false,
      polygon: false,
      polyline: false,
      circle: false,
      circlemarker: false,
      rectangle: false,
    },
    edit: {
      remove: false,
      edit: false,
    },
  };

  @Input('gateways') gateways: Map<any, any>;
  constructor(
    private msgService: MessagesService,
    private mqttService: MqttService,
  ) {}

  addMarker(lon, lat, gw: Gateway) {
    const msg = `
    <h5><b>${ gw.name }</b></h5>
    <div>ID: ${ gw.id }</div>
    <div>External ID: ${ gw.metadata.external_id }</div>
    `;
    const marker: L.Marker = L.marker(
      [lon, lat], {
        icon: L.icon({
          iconUrl: 'assets/images/marker-icon.png',
          iconSize: [50, 50],
        }),
      },
    ).bindPopup(msg);
    this.markersGroup.addLayer(marker);
  }

  refreshCoordinate(gateway: Gateway) {
    this.mqttService.connect({ username: gateway.id, password: gateway.key });
    const topic = 'channels/' + gateway.metadata.ctrl_channel_id + '/messages/req';
    this.mqttService.observe(topic).subscribe((message: IMqttMessage) => {
    });
    this.mqttService.observe(topic).subscribe((message: IMqttMessage) => {
      const long = 43;
      const lat = 54;
      this.addMarker(long, lat, gateway);
    });
  }

  ngOnChanges() {
    this.gateways.forEach((gw) => {
      const channelID: string = gw.metadata ? gw.metadata.data_channel_id : '';

      if (gw.key !== undefined && channelID !== '') {
        this.msgService.getMessages(channelID, gw.key, gw.id).subscribe(
          (resp: any) => {
            let lon: Number;
            let lat: Number;
            if (resp.messages) {
              resp.messages.forEach(msg => {
                // Store lon and lat fields chronologically
                if (msg.name.includes('lon') && !lon) {
                  lon = msg.value;
                }
                if (msg.name.includes('lat') && !lat) {
                  lat = msg.value;
                }
                // Stop for loop if both values are set
                if (lon && lat) {
                  this.addMarker(lon, lat, gw);
                  return;
                }
              });

              this.focusMap();
            }
          },
        );
      }
    });
  }

  onMapReady(map: L.Map) {
    this.map = map;
    this.map.addLayer(this.markersGroup);
    setTimeout(() => {
      map.invalidateSize();
    }, 0);

    map.on('draw:created', e => {
    });
  }

  // Get marker bounds to apply flyToBounds effect on new map view
  focusMap() {
    const markers = this.markersGroup.getLayers();

    if (markers && markers.length) {
      const featureGroup = new L.FeatureGroup(markers);
      const bounds = featureGroup.getBounds();
      const options = {
        maxZoom: 7,
        duration: 2,
      };

      this.map.flyToBounds(bounds, options);
    }
  }
}
