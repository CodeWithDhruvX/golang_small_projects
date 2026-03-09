import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterLink } from '@angular/router';

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.scss']
})
export class SidebarComponent {
  @Input() isOpen = true;
  @Output() close = new EventEmitter<void>();

  navigationItems = [
    {
      name: 'Dashboard',
      href: '/dashboard',
      icon: '📊',
      current: false
    },
    {
      name: 'Documents',
      href: '/documents',
      icon: '📄',
      current: false
    },
    {
      name: 'Upload',
      href: '/upload',
      icon: '📤',
      current: false
    },
    {
      name: 'Chat',
      href: '/chat',
      icon: '💬',
      current: false
    },
    {
      name: 'Search',
      href: '/search',
      icon: '🔍',
      current: false
    }
  ];

  constructor(private router: Router) {}

  onItemClick(item: any): void {
    // Reset current state for all items
    this.navigationItems.forEach(navItem => {
      navItem.current = navItem.href === item.href;
    });
    
    this.close.emit();
  }

  isCurrentRoute(href: string): boolean {
    return this.router.url === href;
  }

  closeSidebar(): void {
    this.close.emit();
  }
}
