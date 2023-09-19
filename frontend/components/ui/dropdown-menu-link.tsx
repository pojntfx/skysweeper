export const DropdownMenuLink = (
  props: React.AnchorHTMLAttributes<HTMLAnchorElement>
) => (
  <a
    role="menuitem"
    className="relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors focus:bg-accent focus:text-accent-foreground hover:bg-accent hover:text-accent-foreground cursor-pointer"
    {...props}
  />
);
