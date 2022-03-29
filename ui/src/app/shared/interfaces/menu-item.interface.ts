export interface MenuItem {
  title?: string;
  icon?: string;
  link?: string;
  home?: boolean;
  children?: MenuItem[];
}
