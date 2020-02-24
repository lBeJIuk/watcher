document.querySelector('#button').addEventListener('click', () => {
  require(['m1/comp1', 'm2/comp3', 'm2/comp4'], function(comp1, comp3, comp4) {
    comp1.log();
    comp3.log();
    comp4.log();
  });
});