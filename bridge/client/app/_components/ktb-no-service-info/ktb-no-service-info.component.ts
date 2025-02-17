import { Component, OnDestroy, OnInit } from '@angular/core';
import { takeUntil } from 'rxjs/operators';
import { DataService } from '../../_services/data.service';
import { Subject } from 'rxjs';

@Component({
  selector: 'ktb-no-service-info',
  templateUrl: './ktb-no-service-info.component.html',
  styleUrls: [],
})
export class KtbNoServiceInfoComponent implements OnInit, OnDestroy {
  private unsubscribe$: Subject<void> = new Subject();
  public isQualityGatesOnly = false;

  constructor(private dataService: DataService) {}

  ngOnInit(): void {
    this.dataService.isQualityGatesOnly.pipe(takeUntil(this.unsubscribe$)).subscribe((isQualityGatesOnly) => {
      this.isQualityGatesOnly = isQualityGatesOnly;
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
