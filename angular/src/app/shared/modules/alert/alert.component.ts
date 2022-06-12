import {
  Component,
  OnInit,
  HostBinding,
  Input,
  ElementRef,
} from '@angular/core';

@Component({
  selector: 'app-alert',
  templateUrl: './alert.component.html',
  styleUrls: ['./alert.component.scss'],
})
export class AlertComponent implements OnInit {
  @Input() title?: string;
  @Input() view: string;
  @Input() removable?: boolean;
  @Input() outline?: boolean;
  @Input() @HostBinding('class.with-before-icon') beforeIcon?: string;
  @Input() @HostBinding('class.with-after-icon') afterIcon?: string;

  @HostBinding('class') get class() {
    return 'alert';
  }
  @HostBinding('class.outline') get getOutline() {
    return this.outline;
  }
  @HostBinding('class.alert-default') get default() {
    return this.view === 'default';
  }
  @HostBinding('class.alert-accent') get accent() {
    return this.view === 'accent';
  }
  @HostBinding('class.alert-success') get success() {
    return this.view === 'success';
  }
  @HostBinding('class.alert-error') get error() {
    return this.view === 'error';
  }
  @HostBinding('class.alert-info') get info() {
    return this.view === 'info';
  }
  @HostBinding('class.alert-warning') get warning() {
    return this.view === 'warning';
  }

  constructor(private element: ElementRef) {
    this.view = 'default';
  }

  ngOnInit() {}

  removeAlert() {
    this.element.nativeElement.remove();
  }
}
