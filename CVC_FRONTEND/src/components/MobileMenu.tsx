import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from './ui/collapsible';
import { Separator } from './ui/separator';
import { ChevronsUpDown } from 'lucide-react';

type MobileMenuProps = {
  navMenu: MenuItem[];
};

import { Button } from './ui/button';

import { MenuItem } from '@/types';

import { NavLink } from 'react-router-dom';
const MobileMenu = ({ navMenu }: MobileMenuProps) => {
  return (
    <div>
      <ul className='mb-3'>
        {navMenu.map(({ href, label, submenu }, index) => (
          <li key={index}>
            {submenu ? (
              <Collapsible>
                <CollapsibleTrigger asChild>
                  <Button
                    variant='ghost'
                    className='w-full justify-between'
                  >
                    {label}
                    <ChevronsUpDown />
                  </Button>
                </CollapsibleTrigger>

                <CollapsibleContent className='ps-2'>
                  <ul className='border-l border-l-muted-foreground/20'>
                    {submenu.map(({ href, label }, index) => (
                      <li key={index}>
                        <Button
                          asChild
                          variant='ghost'
                          className='w-full justify-start text-muted-foreground hover:bg-transparent'
                        >
                          <a href={href}>{label}</a>
                        </Button>
                      </li>
                    ))}
                  </ul>
                </CollapsibleContent>
              </Collapsible>
            ) : (
              <Button
                asChild
                variant='ghost'
                className='w-full p-3 justify-start'
              >
                <a href={href}>{label}</a>
              </Button>
            )}
          </li>
        ))}
      </ul>
      <Separator className='bg-muted-foreground/20' />
      <div className='w-full grid grid-cols-2 gap-2 mt-4'>
     <Button
          variant='ghost'
          className='w-full'
        >
           <NavLink to="/register">Register</NavLink>
        </Button>
        <Button className='w-full '><NavLink to="/login">Login</NavLink></Button>
      </div>
    </div>
  );
};

export default MobileMenu;
