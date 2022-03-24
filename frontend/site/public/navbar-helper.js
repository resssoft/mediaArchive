/* 
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

const SIDEBAR_SMALL_SCREEN_HIDDEN_CLASS = "sidebar-hidden";
const SIDEBAR_LARGE_SCREEN_HIDDEN_CLASS = "sidebar-lg-hidden";

function Navbar(){
    
    this.toggleSidebar = function(){
        let body = document.querySelector('body');
        
        if (body.classList.contains(SIDEBAR_SMALL_SCREEN_HIDDEN_CLASS)) {
            body.classList.remove(SIDEBAR_SMALL_SCREEN_HIDDEN_CLASS);
        }
        else{
            body.classList.add(SIDEBAR_SMALL_SCREEN_HIDDEN_CLASS);
        }
        
        if (body.classList.contains(SIDEBAR_LARGE_SCREEN_HIDDEN_CLASS)) {
            body.classList.remove(SIDEBAR_LARGE_SCREEN_HIDDEN_CLASS);
        }
        else{
            body.classList.add(SIDEBAR_LARGE_SCREEN_HIDDEN_CLASS);
        }
    };
    
}

window.navbar = new Navbar();