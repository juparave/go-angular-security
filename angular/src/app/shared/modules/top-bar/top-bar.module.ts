import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TopBarComponent } from './components/top-bar/top-bar.component';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';

@NgModule({
  declarations: [TopBarComponent],
  imports: [CommonModule, MatToolbarModule, MatButtonModule],
  exports: [TopBarComponent],
})
export class TopBarModule {}
