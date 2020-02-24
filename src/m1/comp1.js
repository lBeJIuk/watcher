define([
  'm1/comp2',
  'm2/comp4'
], function(
  comp2,
  comp4
) {
  return {
    log: function() {
      console.log('log from comp1', '1');
      comp2.log();
      comp4.log();
    }
  }
});
