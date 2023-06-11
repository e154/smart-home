import stream from '@/api/stream';
import {UUID} from 'uuid-generator-ts';

class CustomNavigator {

  private isWatching: boolean = false;

  public watchPosition() {
    if (this.isWatching) {
      return;
    }
    this.isWatching = true;
    if (navigator.geolocation) {
      // navigator.geolocation.getCurrentPosition((position: GeolocationPosition) => console.log(position));
      navigator.geolocation.watchPosition((position: GeolocationPosition) => this.updateLocation(position));
    }
  }

  private updateLocation(position: GeolocationPosition) {

    // console.log('user_id', this.userId);
    // console.log(position);

    stream.send({
      id: UUID.createUUID(),
      query: 'event_update_device_location',
      body: btoa(JSON.stringify({
        lat: position.coords.latitude,
        lon: position.coords.longitude,
        accuracy: position.coords.accuracy,
        speed: position.coords.speed,
      }))
    });
  }
}

const customNavigator = new CustomNavigator();
// customNavigator.watchPosition()

export default customNavigator;
