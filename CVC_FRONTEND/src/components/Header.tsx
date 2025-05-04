
import { Button } from './ui/button';
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover';

import { Menu } from 'lucide-react';
import { navMenu } from '@/constants';
import MobileMenu from './MobileMenu';
import Logo from './Logo';

import { NavLink } from 'react-router-dom';

const Header = () => {
  return (
    <header className="h-16 md:h-20 lg:h-24 flex items-center">
    <div className="container flex justify-between items-center">
      {/* Logo on the left */}
      <div className="flex-shrink-0">
        <Logo variant="icon" />
      </div>

      {/* Spacer to push buttons to the right */}
      <div className="flex-grow"></div>

      {/* Register/Login buttons on the right (hidden on mobile) */}
      <div className="flex items-center gap-2 max-lg:hidden">
        <Button variant="ghost">
          <NavLink to="/register">Register</NavLink>
        </Button>
        <Button>
          <NavLink to="/login">Login</NavLink>
        </Button>
      </div>

      {/* Mobile menu button (visible on mobile) */}
      <Popover>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            size="icon"
            className="lg:hidden"
          >
            <Menu />
          </Button>
        </PopoverTrigger>
        <PopoverContent
          align="end"
          className="bg-background/50 backdrop-blur-3xl border-foreground/5 border-x-0 border-b-0 rounded-lg overflow-hidden"
        >
          <MobileMenu navMenu={navMenu} />
        </PopoverContent>
      </Popover>
    </div>
  </header>
  );
};

export default Header;
