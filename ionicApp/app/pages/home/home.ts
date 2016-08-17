import {Component} from '@angular/core';
import {NavController} from 'ionic-angular';

@Component({
	templateUrl: 'build/pages/home/home.html'
})

export class HomePage {
	constructor(private navCtrl: NavController) {
  
	}

	valveStatus = "on";
	// this.valveStatus = "on";

	onSwitch() {
		if(this.valveStatus == "off") {
			this.valveStatus = "on";
		} else {
			this.valveStatus = "off";
		}
	}

	switchValve(switchToState: Number) {
		var headers = new Headers();
		headers.append('Content-Type', 'application/x-www-form-urlencoded')

		http.post("http://waterapp.guywmoore.com/valve/1", {headers: headers})
			.map(res => res.json())
			.subscribe(
				data => this.valveStatus = data,
				errors => this.logError(errors),
				() => console.log('Valve changed')
			);
	}

	logError(err) {
		console.error('There was an error: ' + err)
	}
}
