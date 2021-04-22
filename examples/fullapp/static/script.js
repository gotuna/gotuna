document.addEventListener('DOMContentLoaded', () => {

  // toggle "navbar-burger"
  const $navbarBurgers = Array.prototype.slice.call(document.querySelectorAll('.navbar-burger'), 0)
  if ($navbarBurgers.length > 0) {
    $navbarBurgers.forEach( el => {
      el.addEventListener('click', () => {
        const target = el.dataset.target
        const $target = document.getElementById(target)
        el.classList.toggle('is-active')
        $target.classList.toggle('is-active')
      })
    })
  }


  // notifications - delete
  (document.querySelectorAll('.notification .delete') || []).forEach(($delete) => {
    const $notification = $delete.parentNode;

    $delete.addEventListener('click', () => {
      $notification.parentNode.removeChild($notification);
    });
  });

  // notifications - auto-hide
  setTimeout(() => {
    const $notifications = Array.prototype.slice.call(document.querySelectorAll('.notification.auto-close'), 0);
    if ($notifications.length > 0) {
      $notifications.forEach( el => {
        el.style.display = 'none'
      })
    }
  }, 3000)

})
