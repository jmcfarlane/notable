/**
 * @fileoverview Describes a bootstrap tab
 */
define([
  'text!templates/note.html',
  'text!templates/tab.html',
  'backbone',
  'underscore',
  'lib/jquery.hotkeys',
  'ace'
],
function(noteDetailTemplate, tabTemplate) {
  return Backbone.View.extend({

    initialize: function() {
      this._editor = null;
      this._tab = null;
    },

    render: function(collection) {
      var note = this.model.toJSON(),
        tabs = this.options.tabs,
        tabContent = this.options.tabContent;

      // No tabs are active
      tabs.find('.active').removeClass('active');

      // Add new tab content
      tabContent.append(_.template(noteDetailTemplate, {
        note: note
      }));

      // Once the tab content is populated, set el.
      this.el = $(tabContent).find('.editor').last().parent();
      this.$el = $(this.el);

      // Chuck the note content into the editor element
      this.$el.find('.editor')
        .text(note.content)
        .height(document.documentElement.clientHeight - 155);

      // Attach the ace editor to the editor element
      this._editor = ace.edit(_.first(this.$el.find('.editor')));
      this._editor.commands.addCommand({
        name: 'save',
        bindKey: {win: 'Ctrl-S',  mac: 'Command-S'},
        exec: _.bind(function(editor) {
          this.save();
        }, this)
      });
      this._editor.commands.removeCommand('gotoline');
      this._editor.focus();
      window._editor = this._editor; // For web driver
      this.$('input').bind('keydown', 'ctrl+s', _.bind(this.save, this));

      // Somehow this seems to result in the right tab saving, but
      // seemms like a defect waiting to happen.
      $(document).bind('keydown', 'ctrl+s', _.bind(this.save, this));
      $(document).bind('keydown', 'esc', _.bind(this.close, this));
      this.$('.close-button').on('click', _.bind(this.close, this));
      this.$('.delete').on('click', _.bind(this.onDelete, this));
      this.$('.save').on('click', _.bind(this.save, this));
      this.$('.save-close').on('click', _.bind(this.saveAndClose, this));

      // Add a new tab
      tabs.append(_.template(tabTemplate, {
        note: note
      }));

      // Track the tab, show it and set event handlers
      this._tab = this.getTab();
      this._tab.on('shown', _.bind(this.shown, this)).tab('show')
      this.$('.subject input').on('keyup', _.bind(this.onSubjectChange, this));
      return this;
    },

    isVisible: function() {
      return this.$el.is(":visible");
    },

    onDelete: function() {
      // TODO: Wrap this in a confirmation modal
      this.model.destroy();
      this.close();
      return false;
    },

    onSubjectChange: function() {
      this._tab.find('span').text(this.$('.subject input').val());
    },

    saveAndClose: function() {
      this.save(_.bind(this.close, this));
      return false;
    },

    show: function() {
      this._tab.tab('show');
    },

    shown: function() {
      if (!this.model.get('subject')) {
        this.$('.subject input').focus();
      }
    },

    notSaved: function() {
      $('.not-saved').fadeIn().delay(3000).fadeOut();
    },

    saved: function() {
      $('.saved').fadeIn().delay(2000).fadeOut();
    },

    save: function(callback) {
      callback = _.isFunction(callback) ? callback : _.identity;
      if (!this.isVisible()) {
        return;
      }
      this.model.save({
        content: this._editor.getValue(),
        password: this.$el.find('.password input').val(),
        subject: this.$el.find('.subject input').val(),
        tags: this.$el.find('.tags input').val()
      }, {
        error: _.bind(function() {
          this.notSaved();
        }, this),
        success: _.bind(function() {
          this.saved();
          callback();
        }, this)
      });
      return false;
    },

    getTab: function() {
      return this.options.tabs.find(this.selector());
    },

    selector: function() {
      return 'a[href=#'+ this.model.get('uid') +']';
    },

    close: function() {
      // Currently you cannot close a tab unless it's active.  This is
      // to simplify bootstrap crazy, see below.
      if (!this.isVisible()) {
        return;
      }

      // Show the tab to the left of this one and remove this tab
      var tabToDelete = this.options.tabs.find(this.selector());
      tabToDelete.parent().prev().find('a').tab('show');
      tabToDelete.parent().remove();

      // Hiding/removing the tab content seems to _piss off_
      // bootstrap, so sorta hide it by setting it's content to
      // nothing (see next step).
      $(tabToDelete.attr('href')).html('');
      this.trigger('destroy');

      // After bootstrap has had time to calm the hell down, remove
      // the tab body (else crazy "all tab bodies disappear" happens).
      setTimeout(function() {
        $(tabToDelete.attr('href')).remove();
      }, 1000);

      return false;
    }

  });
});
